package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"time"
)

type BatchConsumable[T any] interface {
	Consume(msgs []*sarama.ConsumerMessage, ts []T) error
}

type BatchHandler[T any] struct {
	l         logx.Logger
	batchSize int
	maxWait   time.Duration
	consumer  BatchConsumable[T]
}

func NewBatchHandler[T any](l logx.Logger, consumer BatchConsumable[T], opts ...BatchHandlerOption[T]) *BatchHandler[T] {
	h := &BatchHandler[T]{
		l:         l,
		consumer:  consumer,
		batchSize: 10,
		maxWait:   1 * time.Second,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

type BatchHandlerOption[T any] func(bh *BatchHandler[T])

func WithBatchSize[T any](size int) BatchHandlerOption[T] {
	return func(bh *BatchHandler[T]) {
		bh.batchSize = size
	}
}

func WithMaxWait[T any](maxWait time.Duration) BatchHandlerOption[T] {
	return func(bh *BatchHandler[T]) {
		bh.maxWait = maxWait
	}
}

func (h *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for {
		batch := make([]*sarama.ConsumerMessage, 0, h.batchSize)
		ts := make([]T, 0, h.batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), h.maxWait)
		batchDone := false
		var last *sarama.ConsumerMessage
		for _ = range h.batchSize {
			if batchDone {
				break
			}
			select {
			case <-ctx.Done():
				// 达到一批次的最大等待时间
				batchDone = true
			case msg, ok := <-msgs:
				if !ok {
					// chan关闭
					cancel()
					return nil
				}

				last = msg
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					h.l.WithCtx(ctx).Error("batch handler: failed to unmarshal message",
						logx.Error(err),
						logx.String("topic", msg.Topic),
						logx.Int32("partition", msg.Partition),
						logx.Int64("offset", msg.Offset))
					// TODO 考虑丢死信队列
					continue
				}
				batch = append(batch, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		if len(batch) == 0 {
			continue
		}

		err := h.consumer.Consume(batch, ts)
		if err != nil {
			h.l.Error("batch handler: failed to handle message",
				logx.String("topic", last.Topic),
				logx.Error(err),
			)
			// TODO 重试 还是丢到死信队列
		}
		if last != nil {
			session.MarkMessage(last, "")
		}
	}
}
