package ioc

import (
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

//	func InitLogger() logx.Logger {
//		l, err := zap.NewDevelopment(zap.AddCallerSkip(1))
//		if err != nil {
//			panic(err)
//		}
//		logger := logx.NewZapLogger(l)
//		logx.Register(logger)
//		return logger
//	}
func InitLogger() logx.Logger {

	type Config struct {
		Level    string `mapstructure:"level"`
		Filename string `mapstructure:"filename"`
		MaxSize  int    `mapstructure:"maxsize"`
		MaxAge   int    `mapstructure:"maxage"`
	}
	var config Config
	if err := viper.UnmarshalKey("log", &config); err != nil {
		panic(err)
	}
	writers := []io.Writer{
		os.Stdout,
		&lumberjack.Logger{
			Filename:  config.Filename,
			MaxSize:   config.MaxSize,
			MaxAge:    config.MaxAge,
			LocalTime: true,
			Compress:  false,
		},
	}
	logrus.SetOutput(io.MultiWriter(writers...))

	if level, err := logrus.ParseLevel(config.Level); err == nil {
		logrus.SetLevel(level)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	logrus.AddHook(logx.NewContextHook())
	return logx.NewLogrusLogger()
}
