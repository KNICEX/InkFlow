package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
)

type AsyncService interface {
	SubmitInk(ctx context.Context, ink domain.Ink) error
}
type Service interface {
	ReviewInk(ctx context.Context, ink domain.Ink) (domain.ReviewResult, error)
}
type FailoverService interface {
	RetryFail(ctx context.Context) error
	Create(ctx context.Context, typ domain.ReviewType, evt any, er error) error
}
