package logx

import (
	"context"
	"time"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	// WithCtx returns a new Logger with the given context.
	// highly recommended
	WithCtx(ctx context.Context) Logger
	WithField(field ...Field) Logger
	WithSkip(skip int) Logger
}

type Field struct {
	Key   string
	Value any
}

func Any(key string, value any) Field {
	return Field{Key: key, Value: value}
}

func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int32(key string, value int32) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

func Duration(value time.Duration) Field {
	return Field{Key: "duration", Value: value}
}

func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value}
}
