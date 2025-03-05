package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/pkg/logx"
)

type Handler[T any] struct {
	l  logx.Logger
	fn func(msg *sarama.ConsumerMessage, event T) error
}

func NewHandler[T any](fn func(msg *sarama.ConsumerMessage, event T) error, l logx.Logger) *Handler[T] {
	return &Handler[T]{
		l:  l,
		fn: fn,
	}
}

func (h Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			h.l.Error("failed to unmarshal message",
				logx.Error(err),
				logx.String("topic", msg.Topic),
				logx.Int32("partition", msg.Partition),
				logx.Int64("offset", msg.Offset))
		}
		err = h.fn(msg, t)
		if err != nil {
			h.l.Error("failed to handle message",
				logx.Error(err),
				logx.String("topic", msg.Topic),
				logx.Int32("partition", msg.Partition),
				logx.Int64("offset", msg.Offset))
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
