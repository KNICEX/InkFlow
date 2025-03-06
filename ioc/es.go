package ioc

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
)

func InitEs() *elasticsearch.Client {
	type Config struct {
		Addr  string `mapstructure:"addr"`
		Sniff bool   `mapstructure:"sniff"`
	}
	var cfg Config
	if err := viper.UnmarshalKey("es", &cfg); err != nil {
		panic(err)
	}
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.Addr},
	})
	if err != nil {
		panic(err)
	}
	return client
}
