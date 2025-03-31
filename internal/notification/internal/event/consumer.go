package event

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/notification/internal/service"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/saramax"
)

const (
	notificationGroup = "notification-group"
	topicFollow       = "user-follow"
	topicCommentReply = "comment-reply"
	topicCommentLike  = "comment-like"
	topicInkLike      = "ink-like"
)

type NotificationConsumer struct {
	cli      sarama.Client
	svc      service.NotificationService
	handlers map[string]Handler
	l        logx.Logger
}

func NewNotificationConsumer(cli sarama.Client, svc service.NotificationService, l logx.Logger) *NotificationConsumer {
	return &NotificationConsumer{
		cli:      cli,
		svc:      svc,
		handlers: make(map[string]Handler),
		l:        l,
	}
}

func (c *NotificationConsumer) RegisterHandler(handlers ...Handler) error {
	for _, handler := range handlers {
		if _, ok := c.handlers[handler.Topic()]; ok {
			return fmt.Errorf("%s handler already exists", handler.Topic())
		}
		c.handlers[handler.Topic()] = handler
	}
	return nil
}

func (c *NotificationConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient(notificationGroup, c.cli)
	if err != nil {
		return err
	}
	go func() {
		err = cg.Consume(context.Background(),
			[]string{topicFollow, topicCommentReply, topicCommentLike, topicInkLike},
			saramax.NewRawHandler(c.l, c))
		if err != nil {
			c.l.Warn("notification consumer quit...", logx.Error(err))
		}
	}()
	return nil
}

func (c *NotificationConsumer) Consume(msg *sarama.ConsumerMessage) error {
	topic := msg.Topic
	if handler, ok := c.handlers[topic]; ok {
		return handler.HandleMessage(context.Background(), msg)
	} else {
		c.l.Error("no handler found for topic", logx.String("topic", topic))
		return nil
	}

}
