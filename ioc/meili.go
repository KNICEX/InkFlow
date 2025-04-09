package ioc

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/spf13/viper"
)

func InitMeiliSearch() meilisearch.ServiceManager {
	type Config struct {
		Addr      string `mapstructure:"addr"`
		MasterKey string `mapstructure:"master_key"`
	}
	var cfg Config
	err := viper.UnmarshalKey("meilisearch", &cfg)
	if err != nil {
		panic(err)
	}
	return meilisearch.New(cfg.Addr, meilisearch.WithAPIKey(cfg.MasterKey))
}
