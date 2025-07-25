package ioc

import (
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/notification"
	"github.com/KNICEX/InkFlow/internal/recommend"
	"github.com/KNICEX/InkFlow/internal/review"
	"github.com/KNICEX/InkFlow/internal/search"
	"github.com/KNICEX/InkFlow/pkg/saramax"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `mapstructure:"addrs"`
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSyncProducer(client sarama.Client) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}

func InitConsumers(inkRead *interactive.InkViewConsumer, review *review.Consumer,
	search *search.SyncConsumer, notification *notification.SyncConsumer,
	recommend *recommend.SyncConsumer) []saramax.Consumer {
	return []saramax.Consumer{
		inkRead,
		review,
		search,
		notification,
		recommend,
	}
}
