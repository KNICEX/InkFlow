package logx

import (
	"context"
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	entry *logrus.Entry
	skip  int // 保留你的接口定义中的 skip，用不到可以忽略
}

func NewLogrusLogger() Logger {
	return &logrusLogger{
		entry: logrus.NewEntry(logrus.StandardLogger()),
	}
}

func (l *logrusLogger) WithCtx(ctx context.Context) Logger {
	return &logrusLogger{
		entry: logrus.NewEntry(logrus.StandardLogger()).WithContext(ctx),
		skip:  l.skip,
	}
}

func (l *logrusLogger) WithField(fields ...Field) Logger {
	logrusFields := logrus.Fields{}
	for _, f := range fields {
		logrusFields[f.Key] = f.Value
	}
	return &logrusLogger{
		entry: l.entry.WithFields(logrusFields),
		skip:  l.skip,
	}
}

func (l *logrusLogger) WithSkip(skip int) Logger {
	return &logrusLogger{
		entry: l.entry,
		skip:  skip,
	}
}

// Helper to convert Field to logrus.Fields
func toLogrusFields(fields []Field) logrus.Fields {
	logrusFields := logrus.Fields{}
	for _, f := range fields {
		logrusFields[f.Key] = f.Value
	}
	return logrusFields
}

func (l *logrusLogger) Debug(msg string, fields ...Field) {
	l.entry.WithFields(toLogrusFields(fields)).Debug(msg)
}

func (l *logrusLogger) Info(msg string, fields ...Field) {
	l.entry.WithFields(toLogrusFields(fields)).Info(msg)
}

func (l *logrusLogger) Warn(msg string, fields ...Field) {
	l.entry.WithFields(toLogrusFields(fields)).Warn(msg)
}

func (l *logrusLogger) Error(msg string, fields ...Field) {
	l.entry.WithFields(toLogrusFields(fields)).Error(msg)
}

func (l *logrusLogger) Fatal(msg string, fields ...Field) {
	l.entry.WithFields(toLogrusFields(fields)).Fatal(msg)
}
