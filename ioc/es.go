package ioc

import "github.com/elastic/go-elasticsearch/v8"

func InitEs() *elasticsearch.Client {
	type Config struct {
		Url   string `yaml:"url"`
		Sniff bool   `yaml:"sniff"`
	}
	return nil
}
