package schedule

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/review"
)

type ReviewFailoverActivity struct {
	svc review.FailoverService
}

func NewReviewFailoverActivity(svc review.FailoverService) *ReviewFailoverActivity {
	return &ReviewFailoverActivity{
		svc: svc,
	}
}

func (a *ReviewFailoverActivity) RetryFail(ctx context.Context) error {
	return a.svc.RetryFail(ctx)
}
