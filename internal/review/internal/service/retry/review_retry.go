package retry

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
	"github.com/KNICEX/InkFlow/internal/review/internal/repo"
	"github.com/KNICEX/InkFlow/internal/review/internal/service"
	"github.com/KNICEX/InkFlow/internal/review/internal/service/llm"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"go.temporal.io/sdk/client"
)

const reviewSignal = "review-signal"

type reviewRetryService struct {
	workflowCli client.Client
	repo        repo.ReviewFailRepo
	svc         llm.Service
	log         logx.Logger
	batchSize   int
}

func (s *reviewRetryService) Create(ctx context.Context, evt domain.ReviewEvent, er error) error {
	return s.repo.Create(ctx, evt, er)
}

func NewReviewRetryService(repo repo.ReviewFailRepo, svc llm.Service, log logx.Logger, batchSize int) service.ReviewRetryService {
	return &reviewRetryService{
		repo:      repo,
		svc:       svc,
		log:       log,
		batchSize: batchSize,
	}
}

// RetryOnce 是公开暴露的方法，供定时任务调用
func (s *reviewRetryService) RetryOnce(ctx context.Context) error {
	offset := 0
	for {
		events, err := s.repo.Find(ctx, offset, s.batchSize)
		if err != nil {
			s.log.Error("failed to query failed review events", logx.Error(err))
			return err
		}

		if len(events) == 0 {
			s.log.Info("no failed review events to retry")
			return nil
		}

		var successIds []int64

		for _, evt := range events {
			result, err := s.svc.ReviewInk(ctx, evt.Ink)
			if err != nil {
				s.log.Warn("review failed again", logx.String("workflowId", evt.WorkflowId), logx.Error(err))
				continue
			}
			err = s.workflowCli.SignalWorkflow(ctx, evt.WorkflowId, "", reviewSignal, result)
			if err != nil {
				s.log.Error("failed to retry review", logx.String("workflowId", evt.WorkflowId), logx.Error(err))
				continue
			}

			if id := evt.Id; id > 0 {
				successIds = append(successIds, id)
			}
		}

		if len(successIds) > 0 {
			err := s.repo.Delete(ctx, successIds)
			if err != nil {
				s.log.Error("failed to delete retried events", logx.Error(err))
				return err
			}
			s.log.Info("successfully retried and deleted failed review events")
		}
		offset += s.batchSize
	}
}
