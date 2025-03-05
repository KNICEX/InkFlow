package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
)

type InkService interface {
	SaveDraft(ctx context.Context, ink domain.Ink) (int64, error)  // 保存草稿
	Publish(ctx context.Context, ink domain.Ink) (int64, error)    // 发布
	Withdraw(ctx context.Context, ink domain.Ink) error            // 撤回
	GetLiveInk(ctx context.Context, id int64) (domain.Ink, error)  // 获取公开ink
	GetDraftInk(ctx context.Context, id int64) (domain.Ink, error) // 获取草稿ink
	ListLiveByAuthorId(ctx context.Context, authorId int64, maxId int64, limit int) ([]domain.Ink, error)
	ListDraftByAuthorId(ctx context.Context, authorId int64, maxId int64, limit int) ([]domain.Ink, error)
	ListAllLive(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error)
	ListAllDraft(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error)
}
