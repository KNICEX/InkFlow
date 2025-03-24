package event

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/notification/internal/service"
	"github.com/KNICEX/InkFlow/pkg/logx"
)

const (
	notificationGroup = "notification-group"
	topicFollowEvent  = "user-follow-event"
)

type NotificationConsumer struct {
	cli      sarama.Client
	svc      service.NotificationService
	handlers map[string]Handler
	l        logx.Logger
}

func (c *NotificationConsumer) RegisterHandler(topic string, handler Handler) error {
	if _, ok := c.handlers[topic]; ok {
		return fmt.Errorf("%s handler already exists", topic)
	}
	c.handlers[topic] = handler
	return nil
}

func (c *NotificationConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient(notificationGroup, c.cli)
	if err != nil {
		return err
	}
	go func() {
		err = cg.Consume(context.Background(), []string{topicFollowEvent}, c)
		if err != nil {
			c.l.Warn("follow notification consumer quit...", logx.Error(err))
		}
	}()
	return nil
}

func (c *NotificationConsumer) Consume(msg *sarama.ConsumerMessage) error {
	topic := msg.Topic
	if handler, ok := c.handlers[topic]; ok {
		return handler.HandleMessage(msg)
	} else {
		c.l.Error("no handler found for topic", logx.String("topic", topic))
		return nil
	}

}

func (c *NotificationConsumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *NotificationConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *NotificationConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	return nil
}
