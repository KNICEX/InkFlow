package saramax

import (
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"golang.org/x/sync/semaphore"
)

type RawHandler struct {
	l        logx.Logger
	consumer RawConsumable
	se       semaphore.Weighted
}
type RawConsumable interface {
	Consume(msg *sarama.ConsumerMessage) error
}

func NewRawHandler(l logx.Logger, consumer RawConsumable) *RawHandler {
	return &RawHandler{
		l:        l,
		consumer: consumer,
	}
}

func (h *RawHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {

		err := h.consumer.Consume(msg)
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

func (h *RawHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *RawHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}
