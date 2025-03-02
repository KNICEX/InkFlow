package logx

import "context"

var _ Logger = (*NopLogger)(nil)

type NopLogger struct{}

func (l *NopLogger) WithField(field ...Field) Logger {
	return l
}

func NewNopLogger() *NopLogger {
	return &NopLogger{}
}

func (l *NopLogger) Debug(msg string, fields ...Field) {}

func (l *NopLogger) Info(msg string, fields ...Field) {}

func (l *NopLogger) Warn(msg string, fields ...Field) {}

func (l *NopLogger) Error(msg string, fields ...Field) {}

func (l *NopLogger) Fatal(msg string, fields ...Field) {}

func (l *NopLogger) WithCtx(ctx context.Context) Logger {
	return l
}
