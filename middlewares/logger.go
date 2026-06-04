package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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

	if err := os.MkdirAll(filepath.Dir(logsPath), os.ModePerm); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(logsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		panic(err)
	}

	logger := zerolog.New(logFile).With().Timestamp().Logger()

	return func(ctx *gin.Context) {
		requestBody := make(map[string]any)
		formFiles := []map[string]any{}

		logEvent := logger.Info()
		startedAt := time.Now()

		requestContentType := ctx.GetHeader("Content-Type")

		if strings.HasPrefix(requestContentType, "multipart/form-data") {

			if err := ctx.Request.ParseMultipartForm(32 << 20); err == nil && ctx.Request.MultipartForm != nil {
				for k, v := range ctx.Request.MultipartForm.Value {
					if len(v) == 1 {
						requestBody[k] = v[0]
					} else {
						requestBody[k] = v
					}
				}

				for field, files := range ctx.Request.MultipartForm.File {

					for _, file := range files {
						formFiles = append(formFiles, map[string]any{
							"field":        field,
							"file_name":    file.Filename,
							"file_size":    formatByteSize(file.Size),
							"content_type": file.Header.Get("Content-Type"),
						})
					}
				}

				if len(formFiles) > 0 {
					requestBody["form_files"] = formFiles
				}
			}

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
				queries, _ := url.ParseQuery(string(bodyBytes))
				for k, v := range queries {
					if len(v) == 1 {
						requestBody[k] = v[0]
					} else {
						requestBody[k] = v
					}
				}
			}

		}

		customLogResponseWriter := &LogResponseWriter{
			ResponseWriter: ctx.Writer,
			responseBody:   bytes.NewBufferString(""),
		}
		ctx.Writer = customLogResponseWriter

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
