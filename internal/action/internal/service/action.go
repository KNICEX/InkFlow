package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/action/internal/domain"
	"time"
)

type ActionService interface {
	CreateAction(ctx context.Context, action domain.Action) error
	FindActiveUser(ctx context.Context, uids []int64, lastAction time.Time) ([]domain.User, error)
}
