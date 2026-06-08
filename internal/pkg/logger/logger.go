package logger

import "github.com/rs/zerolog"

type LoggerConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	IsDev      string
}

var Log *zerolog.Logger

func NewLogger(config LoggerConfig) {

}
