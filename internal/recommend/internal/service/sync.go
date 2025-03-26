package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/domain"
)

type SyncService interface {
	InputUser(ctx context.Context, user domain.User) error
	InputInk(ctx context.Context, ink domain.Ink) error
	InputFeedback(ctx context.Context, feedback domain.Feedback) error
	DeleteUser(ctx context.Context, userId int64) error
	HiddenInk(ctx context.Context, inkId int64) error
}

type syncService struct {
}

func NewSyncService() SyncService {
	return &syncService{}
}

func (s *syncService) InputUser(ctx context.Context, user domain.User) error {
	return nil
}

func (s *syncService) InputInk(ctx context.Context, ink domain.Ink) error {
	return nil
}
func (s *syncService) InputFeedback(ctx context.Context, feedback domain.Feedback) error {
	return nil
}

func (s *syncService) DeleteUser(ctx context.Context, userId int64) error {
	return nil
}

func (s *syncService) HiddenInk(ctx context.Context, inkId int64) error {
	return nil
}
