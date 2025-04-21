package failover

import (
	"context"
	"errors"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/review/internal/consts"
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
	"github.com/KNICEX/InkFlow/internal/review/internal/event"
	"github.com/KNICEX/InkFlow/internal/review/internal/repo"
	"github.com/KNICEX/InkFlow/internal/review/internal/service"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"go.temporal.io/sdk/client"
)

var (
	ErrRetryFailTooMany = errors.New("fail to retry review too many times")
)

type ReviewService struct {
	workflowCli client.Client
	repo        repo.ReviewFailRepo
	svc         service.Service
	log         logx.Logger
	batchSize   int
}

func NewReviewService(repo repo.ReviewFailRepo, svc service.Service, log logx.Logger) service.FailoverService {
	return &ReviewService{
		repo:      repo,
		svc:       svc,
		log:       log,
		batchSize: 5,
	}
}

func (s *ReviewService) Create(ctx context.Context, typ domain.ReviewType, evt any, er error) error {
	return s.repo.Create(ctx, domain.FailReview{
		Type:  typ,
		Event: evt,
		Error: er,
	}, er)
}

func (s *ReviewService) RetryFail(ctx context.Context) error {
	offset, errCnt := 0, 0

	s.log.Info("start to retry fail review")
	for {
		fails, err := s.repo.Find(ctx, offset, s.batchSize)
		if err != nil {
			s.log.Error("failed to query failed review events", logx.Error(err))
			return err
		}

		if len(fails) == 0 {
			s.log.Info("no failed review events to retry")
			return nil
		}

		var successIds []int64

		for _, fail := range fails {
			switch fail.Type {
			case domain.ReviewTypeInk:
				err = s.retryInk(ctx, fail.Event.(event.ReviewInkEvent))
				if err != nil {
					errCnt++
				} else {
					successIds = append(successIds, fail.Id)
				}
			default:
				s.log.Error("unsupported review type", logx.String("type", string(fail.Type)))
				err = nil
			}

			if errCnt > 3 {
				// 认为审核服务不可用
				return fmt.Errorf("%w, retry error count: %d", ErrRetryFailTooMany, errCnt)
			}

		}

		if len(successIds) > 0 {
			err = s.repo.Delete(ctx, successIds)
			if err != nil {
				s.log.Error("failed to delete retried events", logx.Error(err))
				return err
			}
			s.log.Info("successfully retried and deleted failed review events")
		}
		offset += s.batchSize
	}
}

func (s *ReviewService) retryInk(ctx context.Context, evt event.ReviewInkEvent) error {
	result, err := s.svc.ReviewInk(ctx, evt.Ink)
	if err != nil {
		s.log.Warn("review failed again", logx.String("workflowId", evt.WorkflowId), logx.Error(err))
		return err
	}
	err = s.workflowCli.SignalWorkflow(ctx, evt.WorkflowId, "", consts.ReviewSignal, result)
	if err != nil {
		s.log.Error("failed to retry review", logx.String("workflowId", evt.WorkflowId), logx.Error(err))
		return err
	}
	return nil
}
