package event

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/saramax"
)

const (
	topicUserCreate    = "user-create"
	topicInkView       = "ink-view"
	topicInkLike       = "ink-like"
	topicInkCancelLike = "ink-cancel-like"

	recommendSyncGroup = "recommend-sync-group"
)

type SyncConsumer struct {
	cli      sarama.Client
	handlers map[string]Handler
	l        logx.Logger
}

func NewSyncConsumer(cli sarama.Client, l logx.Logger) *SyncConsumer {
	return &SyncConsumer{
		cli:      cli,
		handlers: make(map[string]Handler),
		l:        l,
	}
}

func (s *SyncConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient(recommendSyncGroup, s.cli)
	if err != nil {
		return err
	}
	go func() {
		err = cg.Consume(context.Background(), []string{topicUserCreate, topicInkLike, topicInkCancelLike}, saramax.NewRawHandler(s.l, s))
		if err != nil {
			s.l.Warn("recommend sync consumer quit...", logx.Error(err))
		}
	}()
	return nil
}

func (s *SyncConsumer) RegisterHandler(handlers ...Handler) error {
	for _, handler := range handlers {
		if _, ok := s.handlers[handler.Topic()]; ok {
			return fmt.Errorf("%s handler already exists", handler.Topic())
		}
		s.handlers[handler.Topic()] = handler
	}
	return nil
}

func (s *SyncConsumer) Consume(msg *sarama.ConsumerMessage) error {
	topic := msg.Topic
	ctx := context.Background()
	if handler, ok := s.handlers[topic]; ok {
		return handler.HandleMessage(ctx, msg)
	} else {
		s.l.WithCtx(ctx).Error("no handler found for topic", logx.String("topic", topic))
		return nil
	}
}
