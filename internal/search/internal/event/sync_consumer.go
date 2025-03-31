package event

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/search/internal/service"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/saramax"
)

const (
	searchSyncGroup = "search-sync-group"
)

const (
	topicUserCreate   = "user-create"
	topicUserUpdate   = "user-update"
	topicCommentReply = "comment-reply"
)

type SyncConsumer struct {
	svc      service.SyncService
	cli      sarama.Client
	handlers map[string]Handler
	l        logx.Logger
}

func NewSyncConsumer(cli sarama.Client, svc service.SyncService, l logx.Logger) *SyncConsumer {
	return &SyncConsumer{
		cli:      cli,
		svc:      svc,
		handlers: make(map[string]Handler),
		l:        l,
	}
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

func (s *SyncConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient(searchSyncGroup, s.cli)
	if err != nil {
		return err
	}
	go func() {
		err = cg.Consume(context.Background(),
			[]string{topicCommentReply, topicUserCreate, topicUserUpdate},
			saramax.NewRawHandler(s.l, s))
		if err != nil {
			s.l.Warn("search sync consumer quit...", logx.Error(err))
		}
	}()
	return nil
}

func (s *SyncConsumer) Consume(msg *sarama.ConsumerMessage) error {
	topic := msg.Topic
	ctx := context.Background()
	if handler, ok := s.handlers[topic]; ok {
		return handler.HandleMessage(ctx, msg)
	} else {
		s.l.WithCtx(ctx).Error("sync search no matched handler", logx.String("topic", topic))
		return nil
	}
}
