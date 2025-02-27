package logger

import (
	"fmt"
	"os"
	"songlibtest/config"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(C *config.Config) *Logger {
	log := zerolog.New(os.Stdout)
	switch C.LogLevel {
	case "Debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "Info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	return &Logger{logger: &log}
}

func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.log("info", message, args...)
}

func (l *Logger) log(level string, message string, args ...interface{}) {
	if len(args) == 0 {
		if level == "debug" {
			l.logger.Debug().Msg(message)
		} else {
			l.logger.Info().Msg(message)
		}
	} else {
		if level == "debug" {
			l.logger.Debug().Msgf(message, args...)
		} else {
			l.logger.Info().Msgf(message, args...)
		}
	}
}
func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(level, msg.Error(), args...)
	case string:
		l.log(level, msg, args...)
	default:
		l.log(level, fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
