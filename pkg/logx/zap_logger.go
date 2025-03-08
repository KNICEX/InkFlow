package logx

import (
	"context"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type ZapLogger struct {
	l   *zap.Logger
	ctx context.Context
}

func NewZapLogger(l *zap.Logger) Logger {
	return &ZapLogger{l: l}
}

func (z *ZapLogger) Debug(msg string, fields ...Field) {
	z.l.Debug(msg, z.toZapFields(fields)...)
}

func (z *ZapLogger) Info(msg string, fields ...Field) {
	z.l.Info(msg, z.toZapFields(fields)...)
}

func (z *ZapLogger) Warn(msg string, fields ...Field) {
	z.l.Warn(msg, z.toZapFields(fields)...)
}

func (z *ZapLogger) Error(msg string, fields ...Field) {
	z.l.Error(msg, z.toZapFields(fields)...)
}

func (z *ZapLogger) Fatal(msg string, fields ...Field) {
	z.l.Fatal(msg, z.toZapFields(fields)...)
}

func (z *ZapLogger) toZapFields(fields []Field) []zap.Field {
	return lo.Map(fields, func(field Field, _ int) zap.Field {
		return zap.Any(field.Key, field.Value)
	})
}

func (z *ZapLogger) WithCtx(ctx context.Context) Logger {
	// TODO 后续整合trace
	return &ZapLogger{l: z.l, ctx: ctx}
}

func (z *ZapLogger) WithField(field ...Field) Logger {
	return &ZapLogger{l: z.l.With(z.toZapFields(field)...)}
}
