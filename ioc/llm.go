package ioc

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"time"
)

func InitGeminiClient() *genai.Client {
	type Config struct {
		Key string `mapstructure:"key"`
	}
	var cfg Config
	if err := viper.UnmarshalKey("llm.gemini", &cfg); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	cli, err := genai.NewClient(ctx, option.WithAPIKey(cfg.Key))
	if err != nil {
		panic(err)
	}
	return cli
}
