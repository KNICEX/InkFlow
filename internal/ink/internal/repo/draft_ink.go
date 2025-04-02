package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo/dao"
	"github.com/KNICEX/InkFlow/pkg/stringx"
	"github.com/samber/lo"
	"strings"
)

var (
	ErrDraftNotFound = dao.ErrDraftNotFound
)

type DraftInkRepo interface {
	Create(ctx context.Context, ink domain.Ink) (int64, error)
	Delete(ctx context.Context, id int64, authorId int64, status ...domain.Status) error
	FindByIdAndAuthorId(ctx context.Context, id, authorId int64, status ...domain.Status) (domain.Ink, error)
	Update(ctx context.Context, ink domain.Ink) error
	UpdateStatus(ctx context.Context, ink domain.Ink) error
	FindByAuthorId(ctx context.Context, authorId int64, offset, limit int, status ...domain.Status) ([]domain.Ink, error)
	FindAll(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error)
}

var _ DraftInkRepo = (*NoCacheDraftInkRepo)(nil)

type NoCacheDraftInkRepo struct {
	dao dao.DraftDAO
}

func NewNoCacheDraftInkRepo(dao dao.DraftDAO) DraftInkRepo {
	return &NoCacheDraftInkRepo{
		dao: dao,
	}
}

func (repo *NoCacheDraftInkRepo) Create(ctx context.Context, ink domain.Ink) (int64, error) {
	id, err := repo.dao.Insert(ctx, repo.domainToEntity(ink))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *NoCacheDraftInkRepo) Delete(ctx context.Context, id int64, authorId int64, status ...domain.Status) error {
	err := repo.dao.Delete(ctx, id, authorId, lo.Map(status, func(item domain.Status, index int) int {
		return item.ToInt()
	})...)
	if err != nil {
		return err
	}
	return nil
}

func (repo *NoCacheDraftInkRepo) FindByIdAndAuthorId(ctx context.Context, id, authorId int64, status ...domain.Status) (domain.Ink, error) {
	ink, err := repo.dao.FindByIdAndAuthorId(ctx, id, authorId, lo.Map(status, func(item domain.Status, index int) int {
		return item.ToInt()
	})...)
	if err != nil {
		return domain.Ink{}, err
	}
	return repo.entityToDomain(ink), nil
}

func (repo *NoCacheDraftInkRepo) Update(ctx context.Context, ink domain.Ink) error {
	err := repo.dao.Update(ctx, repo.domainToEntity(ink))
	if err != nil {
		return err
	}
	return nil
}

func (repo *NoCacheDraftInkRepo) UpdateStatus(ctx context.Context, ink domain.Ink) error {
	err := repo.dao.UpdateStatus(ctx, ink.Id, ink.Author.Id, ink.Status.ToInt())
	if err != nil {
		return err
	}
	return nil
}

func (repo *NoCacheDraftInkRepo) FindByAuthorId(ctx context.Context, authorId int64, offset, limit int, status ...domain.Status) ([]domain.Ink, error) {
	inks, err := repo.dao.FindByAuthorId(ctx, authorId, offset, limit, lo.Map(status, func(item domain.Status, index int) int {
		return item.ToInt()
	})...)
	if err != nil {
		return nil, err
	}
	var result []domain.Ink
	for _, ink := range inks {
		result = append(result, repo.entityToDomain(ink))
	}
	return result, nil
}

func (repo *NoCacheDraftInkRepo) FindAll(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error) {
	inks, err := repo.dao.FindAll(ctx, maxId, limit)
	if err != nil {
		return nil, err
	}
	var result []domain.Ink
	for _, ink := range inks {
		result = append(result, repo.entityToDomain(ink))
	}
	return result, nil
}

func (repo *NoCacheDraftInkRepo) domainToEntity(ink domain.Ink) dao.DraftInk {
	return dao.DraftInk{
		Id:          ink.Id,
		AuthorId:    ink.Author.Id,
		Title:       ink.Title,
		Cover:       ink.Cover,
		Summary:     ink.Summary,
		CategoryId:  ink.Category.Id,
		Tags:        strings.Join(ink.Tags, ","),
		AiTags:      strings.Join(ink.AiTags, ","),
		ContentHtml: ink.ContentHtml,
		ContentMeta: ink.ContentMeta,
		Status:      ink.Status.ToInt(),
		CreatedAt:   ink.CreatedAt,
		UpdatedAt:   ink.UpdatedAt,
	}
}

func (repo *NoCacheDraftInkRepo) entityToDomain(ink dao.DraftInk) domain.Ink {
	return domain.Ink{
		Id: ink.Id,
		Author: domain.Author{
			Id: ink.AuthorId,
		},
		Title:   ink.Title,
		Cover:   ink.Cover,
		Summary: ink.Summary,
		Category: domain.Category{
			Id: ink.CategoryId,
		},
		Tags:        stringx.Split(ink.Tags, ","),
		AiTags:      stringx.Split(ink.AiTags, ","),
		ContentHtml: ink.ContentHtml,
		ContentMeta: ink.ContentMeta,
		Status:      domain.Status(ink.Status),
		CreatedAt:   ink.CreatedAt,
		UpdatedAt:   ink.UpdatedAt,
	}
}
