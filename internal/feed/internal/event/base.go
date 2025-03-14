package event

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/saramax"
)

type FeedEvent struct {
	Type     string
	Metadata map[string]string
}

type FeedEventConsumer struct {
	client sarama.Client
	l      logx.Logger
}

func (f FeedEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("feed", f.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(), []string{"feed"}, saramax.NewBatchHandler[FeedEvent](f.l, f.Consume))
		if err != nil {
			f.l.Warn("consume feed event failed", logx.Error(er))
		}
	}()
	return nil
}

func (f FeedEventConsumer) Consume(msgs []*sarama.ConsumerMessage, events []FeedEvent) error {
	//TODO implement me
	panic("implement me")
}
