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

type DoNothingActionService struct{}

func (s *DoNothingActionService) CreateAction(ctx context.Context, action domain.Action) error {
	return nil
}

func (s *DoNothingActionService) FindActiveUser(ctx context.Context, uids []int64, lastAction time.Time) ([]domain.User, error) {
	res := make([]domain.User, 0, len(uids))
	for _, uid := range uids {
		res = append(res, domain.User{
			Id:           uid,
			LastActionAt: time.Now(),
		})
	}
	return res, nil
}
