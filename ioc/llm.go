package ioc

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"time"
)

func InitGeminiClient() []*genai.Client {
	type Config struct {
		Key []string `mapstructure:"key"`
	}
	var cfg Config
	if err := viper.UnmarshalKey("llm.gemini", &cfg); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	clis := make([]*genai.Client, 0, len(cfg.Key))
	for _, k := range cfg.Key {
		cli, err := genai.NewClient(ctx, option.WithAPIKey(k))
		if err != nil {
			panic(err)
		}
		clis = append(clis, cli)
	}

	return clis
}
