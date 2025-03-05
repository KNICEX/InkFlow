package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/poll/internal/domain"
)

type PollService interface {
	Create(ctx context.Context, vote domain.Poll) (int64, error)
	FindById(ctx context.Context, id int64) (domain.Poll, error)
	Poll(ctx context.Context, uid, vid, optionId int64) error
}
