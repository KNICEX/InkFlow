package ioc

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedisUniversalClient() redis.UniversalClient {
	type Config struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
	}
	var cfg Config
	err := viper.UnmarshalKey("redis", &cfg)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{cfg.Addr},
		Password: cfg.Password,
	})
	if err = redisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return redisClient
}

func InitRedisCmdable(client redis.UniversalClient) redis.Cmdable {
	return client
}
