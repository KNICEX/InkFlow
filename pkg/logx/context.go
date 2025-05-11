package logx

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/sirupsen/logrus"
)

type fieldsKey struct{}

func NewTraceID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func InitContext(ctx context.Context, traceID string) context.Context {
	fields := make(logrus.Fields)
	fields["traceID"] = traceID
	return context.WithValue(ctx, fieldsKey{}, fields)
}

func WithValue(ctx context.Context, field string, value interface{}) {
	fields, ok := ctx.Value(fieldsKey{}).(logrus.Fields)
	if !ok {
		return
	}
	fields[field] = value
}

type ContextHook struct{}

func NewContextHook() *ContextHook {
	return &ContextHook{}
}

// Levels 定义 Hook 适用的日志级别
func (hook *ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels // 适用于所有日志级别
}

// Fire 会在日志输出之前调用，向日志条目中添加从 context 中提取的键值对
func (hook *ContextHook) Fire(entry *logrus.Entry) error {
	if ctx := entry.Context; ctx != nil {
		if fields, ok := ctx.Value(fieldsKey{}).(logrus.Fields); ok {
			for field, value := range fields {
				entry.WithField(field, value)
			}
		}
	}
	return nil
}
