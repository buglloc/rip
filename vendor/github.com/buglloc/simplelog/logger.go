package log

import (
	"os"

	"gopkg.in/inconshreveable/log15.v2"
)

type Logger struct {
	logger log15.Logger
}

func NewLogger() Logger {
	logger := log15.New()
	logger.SetHandler(configFilterHandler(
		log15.StreamHandler(os.Stderr, TextFormat()),
	))
	return Logger{logger: logger}
}

func (l Logger) Debug(msg string, ctx ...interface{}) {
	l.logger.Debug(msg, ctx...)
}

func (l Logger) Info(msg string, ctx ...interface{}) {
	l.logger.Info(msg, ctx...)
}

func (l Logger) Warn(msg string, ctx ...interface{}) {
	l.logger.Warn(msg, ctx...)
}

func (l Logger) Error(msg string, ctx ...interface{}) {
	l.logger.Error(msg, ctx...)
}

func (l Logger) Crit(msg string, ctx ...interface{}) {
	l.logger.Crit(msg, ctx...)
}

func (l Logger) Child(ctx ...interface{}) Logger {
	return Logger{logger: l.logger.New(ctx...)}
}
