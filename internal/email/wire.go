package email

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/email/internal/service"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/spf13/viper"
	"time"
)

type Service = service.Service

func InitService(l logx.Logger) Service {
	type Config struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		FromName string `mapstructure:"from_name"`
	}

	var cfg Config
	if err := viper.UnmarshalKey("email.smtp", &cfg); err != nil {
		panic(err)
	}
	svc, err := service.NewMailService(cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.FromName)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err = svc.Ping(ctx); err != nil {
		l.Error("ping email service error", logx.Error(err))
		//TODO 这个好像是邮件框架的问题，后续处理
	}
	return service.NewAsyncService(svc, l)
}

func InitMemoryService() Service {
	svc := service.NewMemoryService()
	return svc
}
