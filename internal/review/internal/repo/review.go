package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
	"github.com/KNICEX/InkFlow/internal/review/internal/repo/dao"
)

type ReviewFailRepo interface {
	Create(ctx context.Context, evt domain.ReviewEvent, er error) error
	Find(ctx context.Context, offset, limit int) ([]domain.ReviewEvent, error)
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

func (r *reviewFailRepo) Create(ctx context.Context, evt domain.ReviewEvent, er error) error {
	eventJson, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal review event failed: %w", err)
	}

	record := dao.ReviewFail{
		WorkflowId: evt.WorkflowId,
		Event:      string(eventJson),
		Error:      er.Error(),
	}

	return r.dao.Insert(ctx, record)
}

func (r *reviewFailRepo) Find(ctx context.Context, offset, limit int) ([]domain.ReviewEvent, error) {
	records, err := r.dao.Find(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	var events []domain.ReviewEvent
	for _, rec := range records {
		var evt domain.ReviewEvent
		if err := json.Unmarshal([]byte(rec.Event), &evt); err != nil {
			continue
		}
		events = append(events, evt)
	}
	return events, nil
}

func (r *reviewFailRepo) Delete(ctx context.Context, ids []int64) error {
	return r.dao.Delete(ctx, ids)
}
