package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	AppEnv     string
}

type PrettyJSONWriter struct {
	Writer io.Writer
}

type contextKey string

const TraceIdKey contextKey = "trace_id"

var Log *zerolog.Logger

func NewApplicationLogger(config LoggerConfig) {
	Log = NewLogger(config)
}

func NewLogger(config LoggerConfig) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	lvl, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	var writer io.Writer

	if config.AppEnv == "local" {
		if strings.Contains(config.Filename, "app.log") {
			writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		} else {
			writer = PrettyJSONWriter{Writer: os.Stdout}
		}
	} else {
		writer = &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}

func (w PrettyJSONWriter) Write(p []byte) (n int, err error) {
	var prettyJSON bytes.Buffer

	err = json.Indent(&prettyJSON, p, "", "  ")
	if err != nil {
		return w.Writer.Write(p)
	}

	return w.Writer.Write(prettyJSON.Bytes())

	// colorized := colorizeJSON(prettyJSON.Bytes())
	// _, err = w.Writer.Write(colorized)
	// return len(p), err
}

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
	colorGray    = "\033[90m"
	colorBoldRed = "\033[1;31m"
)

func colorizeJSON(p []byte) []byte {
	lines := bytes.Split(p, []byte("\n"))
	var result bytes.Buffer

	for i, line := range lines {
		if len(line) == 0 {
			if i < len(lines)-1 {
				result.WriteByte('\n')
			}
			continue
		}

		q1 := bytes.IndexByte(line, '"')
		if q1 == -1 {
			trimmed := bytes.TrimSpace(line)
			if bytes.Equal(trimmed, []byte("{")) || bytes.Equal(trimmed, []byte("}")) || bytes.Equal(trimmed, []byte("},")) || bytes.Equal(trimmed, []byte("[")) || bytes.Equal(trimmed, []byte("]")) || bytes.Equal(trimmed, []byte("],")) {
				result.WriteString(colorGray)
				result.Write(line)
				result.WriteString(colorReset)
			} else {
				result.Write(line)
			}
		} else {
			q2 := bytes.IndexByte(line[q1+1:], '"')
			if q2 == -1 {
				result.Write(line)
			} else {
				q2 = q1 + 1 + q2
				col := bytes.IndexByte(line[q2+1:], ':')
				if col == -1 {
					result.Write(line)
				} else {
					col = q2 + 1 + col
					mid := bytes.TrimSpace(line[q2+1 : col])
					if len(mid) > 0 {
						result.Write(line)
					} else {
						indent := line[:q1]
						key := line[q1+1 : q2]
						valPart := line[col+1:]

						result.Write(indent)

						keyStr := string(key)
						keyColor := colorCyan
						if keyStr == "error" || keyStr == "panic" || keyStr == "stack" || keyStr == "statck" {
							keyColor = colorRed
						}
						result.WriteString(keyColor)
						result.WriteByte('"')
						result.Write(key)
						result.WriteByte('"')
						result.WriteString(colorReset)

						result.WriteString(colorGray + ":" + colorReset)

						trimmedVal := bytes.TrimSpace(valPart)
						hasComma := bytes.HasSuffix(trimmedVal, []byte(","))
						valWithoutComma := trimmedVal
						if hasComma {
							valWithoutComma = bytes.TrimSuffix(trimmedVal, []byte(","))
						}

						valSpaceLen := len(valPart) - len(trimmedVal)
						if valSpaceLen > 0 {
							result.Write(valPart[:valSpaceLen])
						}

						valStr := string(valWithoutComma)
						valColor := colorWhite

						if keyStr == "level" {
							valStrClean := strings.Trim(valStr, `"`)
							switch valStrClean {
							case "info", "INFO":
								valColor = colorGreen
							case "warn", "warning", "WARN", "WARNING":
								valColor = colorYellow
							case "error", "ERROR":
								valColor = colorBoldRed
							case "panic", "fatal", "PANIC", "FATAL":
								valColor = colorBoldRed
							default:
								valColor = colorGreen
							}
						} else if keyStr == "time" {
							valColor = colorGray
						} else if keyStr == "method" {
							valStrClean := strings.Trim(valStr, `"`)
							switch valStrClean {
							case "GET":
								valColor = colorGreen
							case "POST":
								valColor = colorBlue
							case "PUT":
								valColor = colorYellow
							case "DELETE":
								valColor = colorRed
							default:
								valColor = colorGreen
							}
						} else if keyStr == "path" {
							valColor = colorCyan
						} else if keyStr == "client_ip" {
							valColor = colorMagenta
						} else if keyStr == "panic" {
							valColor = colorBoldRed
						} else if keyStr == "error" {
							valColor = colorRed
						} else if keyStr == "statck_at" || keyStr == "stack_at" {
							valColor = colorYellow
						} else if keyStr == "statck" || keyStr == "stack" {
							valColor = colorGray
						} else {
							if strings.HasPrefix(valStr, `"`) && strings.HasSuffix(valStr, `"`) {
								valColor = colorGreen
							} else if valStr == "true" || valStr == "false" {
								valColor = colorYellow
							} else if valStr == "null" {
								valColor = colorGray
							} else {
								valColor = colorYellow
							}
						}

						result.WriteString(valColor)
						result.Write(valWithoutComma)
						result.WriteString(colorReset)

						if hasComma {
							result.WriteString(colorGray + "," + colorReset)
						}
					}
				}
			}
		}

		if i < len(lines)-1 {
			result.WriteByte('\n')
		}
	}

	return result.Bytes()
}

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIdKey).(string); ok {
		return traceID
	}

	return ""
}
