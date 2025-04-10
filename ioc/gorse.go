package ioc

import (
	"github.com/KNICEX/InkFlow/pkg/gorsex"
	"github.com/spf13/viper"
)

func InitGorseCli() *gorsex.Client {
	type Config struct {
		Addr   string `mapstructure:"addr"`
		ApiKey string `mapstructure:"api_key"`
	}
	var cfg Config
	if err := viper.UnmarshalKey("gorse", &cfg); err != nil {
		panic(err)
	}
	return gorsex.NewClient(cfg.Addr, cfg.ApiKey)
}
