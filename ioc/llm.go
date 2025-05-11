package ioc

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"net/http"
	"net/url"
	"time"
)

func InitGeminiClient() []*genai.Client {
	type Config struct {
		Key []string `mapstructure:"key"`
	}
	type ProxyConfig struct {
		Addr string `mapstructure:"addr"`
	}

	var cfg Config
	var proxyCfg ProxyConfig

	if err := viper.UnmarshalKey("llm.gemini", &cfg); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("proxy", &proxyCfg); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	proxyURL := proxyCfg.Addr
	proxyTransport := &http.Transport{
		Proxy: http.ProxyURL(func() *url.URL {
			u, err := url.Parse(proxyURL)
			if err != nil {
				panic(err)
			}
			return u
		}()),
	}

	// 初始化客户端列表
	clis := make([]*genai.Client, 0, len(cfg.Key))
	for _, k := range cfg.Key {
		cli, err := genai.NewClient(ctx, option.WithAPIKey(k), option.WithHTTPClient(&http.Client{Transport: proxyTransport}))
		if err != nil {
			panic(err)
		}
		clis = append(clis, cli)
	}

	return clis
}
