package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/search/internal/domain"
	"github.com/KNICEX/InkFlow/internal/search/internal/repo/dao"
	"github.com/samber/lo"
)

type InkRepo interface {
	SearchInk(ctx context.Context, query string, offset, limit int) ([]domain.Ink, error)
	InputInk(ctx context.Context, inks []domain.Ink) error
	DeleteInk(ctx context.Context, inkIds []int64) error
}

type inkRepo struct {
	userParser
	dao     dao.InkDAO
	userDAO dao.UserDAO
}

func NewInkRepo(dao dao.InkDAO) InkRepo {
	return &inkRepo{
		dao: dao,
	}
}

func (repo *inkRepo) SearchInk(ctx context.Context, query string, offset, limit int) ([]domain.Ink, error) {
	inkList, err := repo.dao.Search(ctx, query, offset, limit)
	if err != nil || len(inkList) == 0 {
		return nil, err
	}
	authorIds := lo.Map(inkList, func(item dao.Ink, index int) int64 {
		return item.AuthorId
	})
	authors, err := repo.userDAO.SearchByIds(ctx, authorIds)
	return lo.Map(inkList, func(item dao.Ink, index int) domain.Ink {
		ink := repo.entityToDomain(item)
		if author, ok := authors[item.AuthorId]; ok {
			ink.Author = repo.userParser.entityToDomain(author)
		}
		return ink
	}), nil
}

func (repo *inkRepo) InputInk(ctx context.Context, inks []domain.Ink) error {
	err := repo.dao.InputInk(ctx, lo.Map(inks, func(item domain.Ink, index int) dao.Ink {
		return repo.domainToEntity(item)
	}))
	if err != nil {
		return err
	}
	return nil
}

func (repo *inkRepo) DeleteInk(ctx context.Context, inkIds []int64) error {
	return repo.dao.DeleteInk(ctx, inkIds)
}

func (repo *inkRepo) domainToEntity(ink domain.Ink) dao.Ink {
	return dao.Ink{
		Id:        ink.Id,
		Title:     ink.Title,
		AuthorId:  ink.Author.Id,
		Cover:     ink.Cover,
		Content:   ink.Content,
		Tags:      ink.Tags,
		AiTags:    ink.AiTags,
		CreatedAt: ink.CreatedAt,
		UpdatedAt: ink.UpdatedAt,
	}
}

func (repo *inkRepo) entityToDomain(ink dao.Ink) domain.Ink {
	return domain.Ink{
		Id:    ink.Id,
		Title: ink.Title,
		Author: domain.User{
			Id: ink.AuthorId,
		},
		Cover:     ink.Cover,
		Content:   ink.Content,
		Tags:      ink.Tags,
		AiTags:    ink.AiTags,
		CreatedAt: ink.CreatedAt,
		UpdatedAt: ink.UpdatedAt,
	}
}
