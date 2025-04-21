package llm

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
	"github.com/KNICEX/InkFlow/internal/review/internal/event"
	"github.com/KNICEX/InkFlow/internal/review/internal/service"
	"go.temporal.io/sdk/activity"
)

type AsyncWorkflowService struct {
	producer event.ReviewProducer
}

func NewAsyncWorkflowService(producer event.ReviewProducer) service.AsyncService {
	return &AsyncWorkflowService{
		producer: producer,
	}
}

func (a *AsyncWorkflowService) SubmitInk(ctx context.Context, ink domain.Ink) error {
	return a.producer.Produce(ctx, event.ReviewInkEvent{
		Ink:        ink,
		WorkflowId: activity.GetInfo(ctx).WorkflowExecution.ID,
	})
}
