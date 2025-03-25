package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/search/internal/domain"
	"github.com/KNICEX/InkFlow/internal/search/internal/repo"
)

type SyncService interface {
	InputUser(ctx context.Context, users []domain.User) error
	InputInk(ctx context.Context, inks []domain.Ink) error
	InputComment(ctx context.Context, comments []domain.Comment) error
	DeleteInk(ctx context.Context, inkId int64) error
	DeleteUser(ctx context.Context, userId int64) error
	DeleteComment(ctx context.Context, commentId int64) error
}

type syncService struct {
	userRepo    repo.UserRepo
	inkRepo     repo.InkRepo
	commentRepo repo.CommentRepo
}

func NewSyncService(userRepo repo.UserRepo, inkRepo repo.InkRepo, commentRepo repo.CommentRepo) SyncService {
	return &syncService{
		userRepo:    userRepo,
		inkRepo:     inkRepo,
		commentRepo: commentRepo,
	}
}

func (s *syncService) InputUser(ctx context.Context, users []domain.User) error {
	return s.userRepo.InputUser(ctx, users)
}

func (s *syncService) InputInk(ctx context.Context, inks []domain.Ink) error {
	return s.inkRepo.InputInk(ctx, inks)
}

func (s *syncService) InputComment(ctx context.Context, comments []domain.Comment) error {
	return s.commentRepo.InputComment(ctx, comments)
}

func (s *syncService) DeleteInk(ctx context.Context, inkId int64) error {
	err := s.inkRepo.DeleteInk(ctx, []int64{inkId})
	if err != nil {
		return err
	}
	return s.commentRepo.DeleteByBiz(ctx, domain.BizInk, inkId)
}

func (s *syncService) DeleteUser(ctx context.Context, userId int64) error {
	return s.userRepo.DeleteUser(ctx, []int64{userId})
}

func (s *syncService) DeleteComment(ctx context.Context, commentId int64) error {
	return s.commentRepo.DeleteComment(ctx, commentId)
}
