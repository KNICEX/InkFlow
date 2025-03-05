package ioc

import (
	"github.com/KNICEX/InkFlow/pkg/logx"
	"go.uber.org/zap"
)

func InitLogger() logx.Logger {
	l, err := zap.NewDevelopment(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger := logx.NewZapLogger(l)
	logx.Register(logger)
	return logger
}
