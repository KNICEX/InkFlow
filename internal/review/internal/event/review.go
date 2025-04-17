package event

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
	"github.com/KNICEX/InkFlow/internal/review/internal/service"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/saramax"
	"go.temporal.io/sdk/client"
)

const inkReviewTopic = "ink-review"
const inkReviewGroup = "ink-review-group"
const reviewSignal = "review-signal"

type ReviewEvent struct {
	WorkflowId string
	Ink        domain.Ink
}

type ReviewProducer interface {
	Produce(ctx context.Context, event ReviewEvent) error
}

type KafkaReviewProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaReviewProducer(producer sarama.SyncProducer) ReviewProducer {
	return &KafkaReviewProducer{
		producer: producer,
	}
}

func (p *KafkaReviewProducer) Produce(ctx context.Context, event ReviewEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: inkReviewTopic,
		Value: sarama.ByteEncoder(data),
	})
	return err
}

type ReviewConsumer struct {
	workflowCli client.Client
	svc         service.Service
	retrySvc    service.ReviewRetryService
	saramaCli   sarama.Client
	l           logx.Logger
}

func NewReviewConsumer(workflowCli client.Client, svc service.Service, saramaCli sarama.Client, l logx.Logger) *ReviewConsumer {
	return &ReviewConsumer{
		workflowCli: workflowCli,
		svc:         svc,
		saramaCli:   saramaCli,
		l:           l,
	}
}

func (c *ReviewConsumer) Start() error {
	group, err := sarama.NewConsumerGroupFromClient(inkReviewGroup, c.saramaCli)
	if err != nil {
		return err
	}
	go func() {
		err = group.Consume(context.Background(), []string{inkReviewTopic}, saramax.NewHandler(c, c.l))
		if err != nil {
			c.l.Warn("ink review consumer quit...", logx.Error(err))
		}
	}()
	return nil
}

func (c *ReviewConsumer) Consume(msg *sarama.ConsumerMessage, event ReviewEvent) error {
	ctx := context.Background()
	var result domain.ReviewResult
	var err error
	maxRetry := 3

	for i := 0; i < maxRetry; i++ {
		result, err = c.svc.ReviewInk(ctx, event.Ink)
		if err == nil {
			break
		}
		c.l.Warn("review ink failed, will retry", logx.Any("retry", i+1), logx.Error(err))
	}

	if err != nil {
		c.l.Error("review ink failed after retries, storing to fallback queue",
			logx.Error(err), logx.String("workflowId", event.WorkflowId), logx.Any("ink", event.Ink))

		saveErr := c.retrySvc.Create(ctx, domain.ReviewEvent{
			Ink:        event.Ink,
			WorkflowId: event.WorkflowId,
		}, err)
		if saveErr != nil {
			c.l.Error("failed to store failed review", logx.Error(saveErr))
		}
		return err
	}
	return c.workflowCli.SignalWorkflow(ctx, event.WorkflowId, "", reviewSignal, result)
}
