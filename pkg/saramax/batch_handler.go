package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"time"
)

type batchHandleFunc[T any] = func(msgs []*sarama.ConsumerMessage, ts []T) error

type BatchHandler[T any] struct {
	l         logx.Logger
	batchSize int
	maxWait   time.Duration
	fn        batchHandleFunc[T]
}

func NewBatchHandler[T any](l logx.Logger, fn batchHandleFunc[T], opts ...BatchHandlerOption[T]) *BatchHandler[T] {
	h := &BatchHandler[T]{
		l:         l,
		fn:        fn,
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
		for _ = range h.batchSize {
			select {
			case <-ctx.Done():
			case msg, ok := <-msgs:
				if !ok {
					// chan关闭
					cancel()
					return nil
				}

				batch = append(batch, msg)
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					h.l.Error("batch handler: failed to unmarshal message",
						logx.Error(err),
						logx.String("topic", msg.Topic),
						logx.Int32("partition", msg.Partition),
						logx.Int64("offset", msg.Offset))
					continue
				}
				ts = append(ts, t)
			}
		}
		cancel()

		err := h.fn(batch, ts)
		if err != nil {
			h.l.Error("batch handler: failed to handle message",
				logx.Error(err),
			)
			// TODO 重试 还是丢到死信队列
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}
