package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogResponseWriter struct {
	gin.ResponseWriter
	responseBody *bytes.Buffer
}

func (w *LogResponseWriter) Write(data []byte) (int, error) {
	w.responseBody.Write(data)
	return w.ResponseWriter.Write(data)
}

func LoggerMiddleware() gin.HandlerFunc {
	logsPath := "logs/http.log"

	// if err := os.MkdirAll(filepath.Dir(logsPath), os.ModePerm); err != nil {
	// 	panic(err)
	// }

	// logFile, err := os.OpenFile(logsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	// if err != nil {
	// 	panic(err)
	// }

	logger := zerolog.New(&lumberjack.Logger{
		Filename:   logsPath,
		MaxSize:    1, // megabytes
		MaxBackups: 7,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
		LocalTime:  true,
	}).With().Timestamp().Logger()

	return func(ctx *gin.Context) {
		requestBody := make(map[string]any)
		formFiles := []map[string]any{}

		logEvent := logger.Info()
		startedAt := time.Now()

		ctx.Next()

		statusCode := ctx.Writer.Status()
		responseContentType := ctx.Writer.Header().Get("Content-Type")
		rawResBody := customLogResponseWriter.responseBody.String()
		var responseBodyParsed any

		if strings.HasPrefix(responseContentType, "image/") {
		} else if strings.HasPrefix(responseContentType, "application/json") ||
			strings.HasPrefix(strings.TrimSpace(rawResBody), "{") ||
			strings.HasPrefix(strings.TrimSpace(rawResBody), "[") {
			if err := json.Unmarshal([]byte(rawResBody), &responseBodyParsed); err != nil {
				responseBodyParsed = rawResBody
			}
		} else {
			responseBodyParsed = rawResBody
		}

		if statusCode >= 500 {
			logEvent = logger.Error()
		} else if statusCode >= 400 {
			logEvent = logger.Warn()
		} else {
			logEvent = logger.Info()
		}

		logEvent.
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.Path).
			Str("query", ctx.Request.URL.RawQuery).
			Str("client_ip", ctx.ClientIP()).
			Str("user_agent", ctx.Request.UserAgent()).
			Str("referer", ctx.Request.Referer()).
			Str("protocol", ctx.Request.Proto).
			Str("host", ctx.Request.Host).
			Str("remote_addr", ctx.Request.RemoteAddr).
			Str("request_uri", ctx.Request.RequestURI).
			Int64("content_length", ctx.Request.ContentLength).
			Interface("headers", ctx.Request.Header).
			Interface("request_body", requestBody).
			Int("status_code", statusCode).
			Interface("response_body", responseBodyParsed).
			Int64("duration_ms", time.Since(startedAt).Milliseconds()).
			Msg("HTTP Request Log")

	}
}

func formatByteSize(size int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
	)

	switch {
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
