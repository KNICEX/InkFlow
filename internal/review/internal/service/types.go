package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
)

type Service interface {
	ReviewInk(ctx context.Context, ink domain.Ink) (domain.ReviewResult, error)
}

type AsyncService interface {
	SubmitInk(ctx context.Context, ink domain.Ink) error
}
