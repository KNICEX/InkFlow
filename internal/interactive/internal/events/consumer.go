package events

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/domain"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/repo"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/saramax"
	"time"
)

const (
	inkViewConsumerGroup = "ink-view-consumer-group"
)

type InkViewConsumer struct {
	client sarama.Client
	repo   repo.InteractiveRepo
	l      logx.Logger
}

func NewInkViewConsumer(client sarama.Client, repo repo.InteractiveRepo, l logx.Logger) *InkViewConsumer {
	return &InkViewConsumer{
		client: client,
		repo:   repo,
		l:      l,
	}
}

func (c *InkViewConsumer) Consume(msgs []*sarama.ConsumerMessage, ts []InkViewEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	inkIds, uids := make([]int64, len(ts)), make([]int64, len(ts))
	for i, t := range ts {
		inkIds[i] = t.InkId
		uids[i] = t.UserId
	}
	return c.repo.IncrViewBatch(ctx, domain.BizInk, inkIds, uids)
}

func (c *InkViewConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient(inkViewConsumerGroup, c.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{inkViewTopic}, saramax.NewBatchHandler(c.l, c))
		if er != nil {
			c.l.Warn("consume quited", logx.Error(er))
		}
	}()
	return nil
}
