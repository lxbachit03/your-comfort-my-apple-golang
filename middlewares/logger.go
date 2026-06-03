package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func LoggerMiddleware() gin.HandlerFunc {

	logsPath := "logs/http.log"

	if err := os.MkdirAll(filepath.Dir(logsPath), os.ModePerm); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(logsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		panic(err)
	}

	logger := zerolog.New(logFile).With().Timestamp().Logger()

	return func(ctx *gin.Context) {
		logEvent := logger.Info()
		startedAt := time.Now()
		requestBody := make(map[string]any)

		requestContentType := ctx.GetHeader("Content-Type")

		log.Print("requestContentType: ", requestContentType)

		if strings.HasPrefix(requestContentType, "multipart/form-data") {

		} else {
			bodyBytes, err := io.ReadAll(ctx.Request.Body)

			if err != nil {
				logger.Error().Err(err).Msg("failed to read request body")
			}

			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			log.Print("requestBody: ", string(bodyBytes))

			if strings.HasPrefix(requestContentType, "application/json") {
				if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
					logger.Error().Err(err).Msg("failed to unmarshal request body")
				}
			} else if strings.HasPrefix(requestContentType, "application/x-www-form-urlencoded") {
			}

		}

		ctx.Next()

		statusCode := ctx.Writer.Status()

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
			// Interface("response_body", responseBodyParsed).
			Int64("duration_ms", time.Since(startedAt).Milliseconds()).
			Msg("HTTP Request Log")

	}
}
