package log

import (
	bugsnag "github.com/bugsnag/bugsnag-go"
	bugsnagerr "github.com/bugsnag/bugsnag-go/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Debug  = "debug"
	Info   = "info"
	Warn   = "warn"
	Error  = "error"
	DPanic = "dpanic"
	Panic  = "panic"
	Fatal  = "fatal"
)

// Logger is the fundamental interface for all log operations. Logger creates a
// log event from keyvals, a variadic sequence of alternating keys and values.
//
// Implementations must be safe for concurrent use by multiple goroutines. In
// particular, any implementation of Logger that appends to keyvals or modifies
// or retains any of its elements must make a copy first.
type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(err error, keyvals ...interface{})
	Fatal(err error, keyvalls ...interface{})
	Bug(err error, keyvalls ...interface{})
	With(args ...interface{}) Logger
}

type logger struct {
	zap        *zap.SugaredLogger
	production bool
}

func (l *logger) Debug(msg string, keyvals ...interface{}) {
	l.zap.Debugw(msg, keyvals...)
}

func (l *logger) Info(msg string, keyvals ...interface{}) {
	l.zap.Infow(msg, keyvals...)
}

func (l *logger) Warn(msg string, keyvals ...interface{}) {
	l.zap.Warnw(msg, keyvals...)
}

func (l *logger) Error(err error, keyvals ...interface{}) {
	if err == nil {
		return
	}
	l.zap.Errorw(err.Error(), keyvals...)
}

func (l *logger) Fatal(err error, keyvals ...interface{}) {
	l.zap.Fatalw(err.Error(), keyvals...)
}

func (l *logger) Bug(bug error, rawData ...interface{}) {
	if !l.production {
		l.zap.Errorw(bug.Error(), rawData...)
		return
	}
	keyvals := make(map[string]interface{})
	for i := 0; i < len(rawData); {
		k, v := rawData[i], rawData[i+1]
		key, ok := k.(string)
		if !ok {
			l.zap.Errorw("non-string key passed")
			return
		}
		keyvals[key] = v
		i += 2
	}
	data := bugsnag.MetaData{
		"Data": keyvals,
	}
	err := bugsnag.Notify(bugsnagerr.New(bug, 1), data)
	if err != nil {
		l.Error(err)
	}
}

func (l *logger) With(args ...interface{}) Logger {
	return &logger{
		production: l.production,
		zap:        l.zap.With(args...),
	}

}

func NewLogger(production bool, level string) (Logger, error) {
	var zapLogger *zap.Logger
	var config zap.Config

	// parse log level
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}

	// use JSON encoder for structured logging in production
	if production {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}
	config.DisableStacktrace = true
	config.Level.SetLevel(logLevel)
	zapLogger, err = config.Build()

	if err != nil {
		return nil, err
	}
	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(1))
	return &logger{
		zap:        zapLogger.Sugar(),
		production: production,
	}, nil
}

func NopLogger() Logger {
	return &logger{zap: zap.NewNop().Sugar()}
}
