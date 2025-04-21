package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
	"github.com/KNICEX/InkFlow/internal/review/internal/repo/dao"
)

var (
	ErrUnknownReviewType = errors.New("unknown review type")
)

type ReviewFailRepo interface {
	Create(ctx context.Context, evt domain.FailReview, er error) error
	Find(ctx context.Context, offset, limit int) ([]domain.FailReview, error)
	Delete(ctx context.Context, ids []int64) error
}

type reviewFailRepo struct {
	dao dao.ReviewFailDAO
}

func NewReviewFailRepo(dao dao.ReviewFailDAO) ReviewFailRepo {
	return &reviewFailRepo{
		dao: dao,
	}
}

func (r *reviewFailRepo) Create(ctx context.Context, evt domain.FailReview, er error) error {
	eventJson, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal review event failed: %w", err)
	}

	record := dao.ReviewFail{
		Event: string(eventJson),
		Error: er.Error(),
	}

	return r.dao.Insert(ctx, record)
}

func (r *reviewFailRepo) Find(ctx context.Context, offset, limit int) ([]domain.FailReview, error) {
	records, err := r.dao.Find(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	res := make([]domain.FailReview, 0, len(records))
	for _, rec := range records {
		domainRec, err := r.toDomain(&rec)
		if err != nil {
			return nil, err
		}
		res = append(res, domainRec)
	}

	return res, nil
}

func (r *reviewFailRepo) Delete(ctx context.Context, ids []int64) error {
	return r.dao.Delete(ctx, ids)
}

func (r *reviewFailRepo) toDomain(entity *dao.ReviewFail) (domain.FailReview, error) {
	var evt any

	switch domain.ReviewType(entity.Type) {
	case domain.ReviewTypeInk:
		if err := json.Unmarshal([]byte(entity.Event), &evt); err != nil {
			return domain.FailReview{}, err
		}
	default:
		return domain.FailReview{}, fmt.Errorf("%w : %s", ErrUnknownReviewType, entity.Type)
	}
	return domain.FailReview{
		Id:        entity.Id,
		Type:      domain.ReviewType(entity.Type),
		Event:     evt,
		Error:     errors.New(entity.Error),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}
