package log

import (
	"os"

	"gopkg.in/inconshreveable/log15.v2"
)

var logger log15.Logger
var maxLvl = InfoLevel

const (
	CritLevel log15.Lvl = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

func init() {
	logger = log15.New()
	logger.SetHandler(configFilterHandler(
		log15.StreamHandler(os.Stderr, TextFormat()),
	))
}

func configFilterHandler(h log15.Handler) log15.Handler {
	return log15.FilterHandler(func(r *log15.Record) (pass bool) {
		return r.Lvl <= maxLvl
	}, h)
}

func SetLevel(level log15.Lvl) {
	maxLvl = level
}

func Debug(msg string, ctx ...interface{}) {
	logger.Debug(msg, ctx...)
}

func Info(msg string, ctx ...interface{}) {
	logger.Info(msg, ctx...)
}

func Warn(msg string, ctx ...interface{}) {
	logger.Warn(msg, ctx...)
}

func Error(msg string, ctx ...interface{}) {
	logger.Error(msg, ctx...)
}

func Crit(msg string, ctx ...interface{}) {
	logger.Crit(msg, ctx...)
}
