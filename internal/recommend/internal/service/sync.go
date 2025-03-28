package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/domain"
)

type SyncService interface {
	InputUser(ctx context.Context, user domain.User) error
	InputInk(ctx context.Context, ink domain.Ink) error
	InputFeedback(ctx context.Context, feedback domain.Feedback) error
	InputRelation(ctx context.Context, relation domain.Relation) error
	DeleteUser(ctx context.Context, userId int64) error
	HiddenInk(ctx context.Context, inkId int64) error
}
