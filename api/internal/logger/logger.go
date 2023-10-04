package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	zerolog.Logger
}

func NewLogger() *Logger {
	return &Logger{Logger: log.Logger}
}


func (l *Logger) Info(msg string, fields ...interface{}) {
	l.Logger.Info().Fields(fields...).Msg(msg)
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	l.Logger.Error().Fields(fields...).Msg(msg)
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.Logger.Debug().Fields(fields...).Msg(msg)
}

func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.Logger.Fatal().Fields(fields...).Msg(msg)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.Logger.Warn().Fields(fields...).Msg(msg)
}

func (l *Logger) Panic(msg string, fields ...interface{}) {
	l.Logger.Panic().Fields(fields...).Msg(msg)
}

