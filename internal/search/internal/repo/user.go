package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/search/internal/domain"
	"github.com/KNICEX/InkFlow/internal/search/internal/repo/dao"
	"github.com/samber/lo"
)

type UserRepo interface {
	Search(ctx context.Context, query string, offset, limit int) ([]domain.User, error)
	InputUser(ctx context.Context, users []domain.User) error
	DeleteUser(ctx context.Context, userId []int64) error
}

type userRepo struct {
	userParser
	dao dao.UserDAO
}

func NewUserRepo(dao dao.UserDAO) UserRepo {
	return &userRepo{
		dao: dao,
	}
}
func (repo *userRepo) Search(ctx context.Context, query string, offset, limit int) ([]domain.User, error) {
	userList, err := repo.dao.Search(ctx, query, offset, limit)
	if err != nil || len(userList) == 0 {
		return nil, err
	}
	return lo.Map(userList, func(item dao.User, index int) domain.User {
		return repo.entityToDomain(item)
	}), nil
}

func (repo *userRepo) InputUser(ctx context.Context, users []domain.User) error {
	err := repo.dao.InputUser(ctx, lo.Map(users, func(item domain.User, index int) dao.User {
		return repo.domainToEntity(item)
	}))
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepo) DeleteUser(ctx context.Context, userIds []int64) error {
	return repo.dao.DeleteUser(ctx, userIds)
}

type userParser struct {
}

func (p userParser) domainToEntity(user domain.User) dao.User {
	return dao.User{
		Id:        user.Id,
		Username:  user.Username,
		Account:   user.Account,
		Avatar:    user.Avatar,
		AboutMe:   user.AboutMe,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (p userParser) entityToDomain(user dao.User) domain.User {
	return domain.User{
		Id:        user.Id,
		Username:  user.Username,
		Account:   user.Account,
		Avatar:    user.Avatar,
		AboutMe:   user.AboutMe,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
