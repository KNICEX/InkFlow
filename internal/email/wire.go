package email

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/email/internal/service"
	"github.com/spf13/viper"
	"time"
)

type Service = service.Service

func InitService() Service {
	type Config struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		FromName string `yaml:"from_ame"`
	}

	var cfg Config
	if err := viper.UnmarshalKey("email.smtp", &cfg); err != nil {
		panic(err)
	}
	svc, err := service.NewSmtpService(cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.FromName)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err = svc.Ping(ctx); err != nil {
		panic(err)
	}
	return svc
}
