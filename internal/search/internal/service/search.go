package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/search/internal/domain"
	"github.com/KNICEX/InkFlow/internal/search/internal/repo"
)

type SearchService interface {
	SearchUser(ctx context.Context, query string, offset, limit int) ([]domain.User, error)
	SearchInk(ctx context.Context, query string, offset, limit int) ([]domain.Ink, error)
	SearchComment(ctx context.Context, query string, offset, limit int) ([]domain.Comment, error)
}

type searchService struct {
	userRepo    repo.UserRepo
	inkRepo     repo.InkRepo
	commentRepo repo.CommentRepo
}

func NewSearchService(userRepo repo.UserRepo, inkRepo repo.InkRepo, commentRepo repo.CommentRepo) SearchService {
	return &searchService{
		userRepo:    userRepo,
		inkRepo:     inkRepo,
		commentRepo: commentRepo,
	}
}

func (s *searchService) SearchUser(ctx context.Context, query string, offset, limit int) ([]domain.User, error) {
	users, err := s.userRepo.Search(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *searchService) SearchInk(ctx context.Context, query string, offset, limit int) ([]domain.Ink, error) {
	inks, err := s.inkRepo.SearchInk(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	return inks, nil
}

func (s *searchService) SearchComment(ctx context.Context, query string, offset, limit int) ([]domain.Comment, error) {
	comments, err := s.commentRepo.Search(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
