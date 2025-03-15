package service

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo"
	"github.com/KNICEX/InkFlow/pkg/logx"
)

var (
	ErrDraftNotFound = repo.ErrDraftNotFound
	ErrLiveNotFound  = repo.ErrLiveInkNotFound
)

// InkService
// 这里将Save和Publish分开，Sava只当作保存草稿，前端编辑时也会定时调用放丢失
// 直接发布文章，前端先调用Save保存草稿，然后再调用Publish发布
type InkService interface {
	Save(ctx context.Context, ink domain.Ink) (int64, error)                 // 保存草稿
	Publish(ctx context.Context, ink domain.Ink) (int64, error)              // 发布
	Withdraw(ctx context.Context, ink domain.Ink) error                      // 撤回
	GetLiveInk(ctx context.Context, id int64) (domain.Ink, error)            // 获取公开ink
	GetDraftInk(ctx context.Context, authorId, id int64) (domain.Ink, error) // 获取草稿ink
	ListLiveByAuthorId(ctx context.Context, authorId int64, offset int, limit int) ([]domain.Ink, error)
	ListPendingByAuthorId(ctx context.Context, authorId int64, offset int, limit int) ([]domain.Ink, error)
	ListReviewFailedByAuthorId(ctx context.Context, authorId int64, offset int, limit int) ([]domain.Ink, error)
	ListDraftByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error)
	ListAllLive(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error)
	ListAllDraft(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error)
}

type inkService struct {
	liveRepo  repo.LiveInkRepo
	draftRepo repo.DraftInkRepo
	l         logx.Logger
}

func NewInkService(liveRepo repo.LiveInkRepo, draftRepo repo.DraftInkRepo, l logx.Logger) InkService {
	return &inkService{
		liveRepo:  liveRepo,
		draftRepo: draftRepo,
	}
}

func (svc *inkService) Save(ctx context.Context, ink domain.Ink) (int64, error) {
	// 保存草稿，草稿状态一定为未发布
	ink.Status = domain.InkStatusUnPublished
	if ink.Id == 0 {
		return svc.draftRepo.Create(ctx, ink)
	}
	return ink.Id, svc.draftRepo.Update(ctx, ink)
}

// Publish
// 从草稿库查询文章(因为更新文章都是先save更新草稿，我们认为草稿是最新的数据)
func (svc *inkService) Publish(ctx context.Context, ink domain.Ink) (int64, error) {
	draft, err := svc.draftRepo.FindByIdAndAuthorId(ctx, ink.Id, ink.Author.Id)
	if err != nil {
		return 0, err
	}
	if draft.Status != domain.InkStatusUnPublished {
		svc.l.WithCtx(ctx).Error("InkService Publish  draft status error",
			logx.Error(err), logx.Int64("inkId", ink.Id), logx.Any("status", draft.Status))
		return 0, errors.New("invalid status")
	}

	// TODO 这里后面需要添加审核
	go func() {
		// 文章发布后，草稿的状态需要更新为已发布
		if er := svc.draftRepo.UpdateStatus(ctx, draft.Id, draft.Author.Id, domain.InkStatusPublished); er != nil {
			svc.l.WithCtx(ctx).Error("InkService Publish  update draft status error",
				logx.Error(er), logx.Int64("inkId", ink.Id), logx.Any("status", draft.Status))
		}
	}()

	draft.Status = domain.InkStatusPublished
	return svc.liveRepo.Save(ctx, draft)
}

// UpdateInkStatus
// 此方法使用权不应该交给用户，仅内部服务可以调用
func (svc *inkService) UpdateInkStatus(ctx context.Context, id int64, authorId int64, status domain.Status) error {
	err := svc.liveRepo.UpdateStatus(ctx, domain.Ink{
		Id: id,
		Author: domain.Author{
			Id: authorId,
		},
		Status: status,
	})
	return err
}

func (svc *inkService) Withdraw(ctx context.Context, ink domain.Ink) error {
	// TODO 暂时设置为私有，后续考虑要不要添加更多状态
	ink.Status = domain.InkStatusPrivate
	err := svc.liveRepo.UpdateStatus(ctx, ink)
	if err != nil {
		return err
	}
	// 更新草稿状态, 让草稿可以继续编辑
	err = svc.draftRepo.UpdateStatus(ctx, ink.Author.Id, ink.Id, domain.InkStatusUnPublished)
	if err != nil {
		return err
	}
	return nil
}

func (svc *inkService) GetLiveInk(ctx context.Context, id int64) (domain.Ink, error) {
	return svc.liveRepo.FindByIdAndStatus(ctx, id, domain.InkStatusPublished)
}

func (svc *inkService) GetDraftInk(ctx context.Context, authorId, id int64) (domain.Ink, error) {
	return svc.draftRepo.FindByIdAndAuthorId(ctx, authorId, id)
}

func (svc *inkService) ListLiveByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error) {
	return svc.liveRepo.ListByAuthorIdAndStatus(ctx, authorId, domain.InkStatusPublished, offset, limit)
}

func (svc *inkService) ListPendingByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error) {
	return svc.liveRepo.ListByAuthorIdAndStatus(ctx, authorId, domain.InkStatusPending, offset, limit)
}

func (svc *inkService) ListReviewFailedByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error) {
	return svc.liveRepo.ListByAuthorIdAndStatus(ctx, authorId, domain.InkStatusReviewFailed, offset, limit)
}

func (svc *inkService) ListDraftByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error) {
	return svc.draftRepo.ListByAuthorId(ctx, authorId, offset, limit)
}

// ListAllLive 这个接口不应该让用户直接使用
func (svc *inkService) ListAllLive(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error) {
	return svc.liveRepo.FindAllByStatus(ctx, domain.InkStatusPublished, maxId, limit)
}

func (svc *inkService) ListAllDraft(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error) {
	return svc.draftRepo.ListAll(ctx, maxId, limit)
}

func (svc *inkService) ListAllReviewFailed(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error) {
	return svc.liveRepo.FindAllByStatus(ctx, domain.InkStatusReviewFailed, maxId, limit)
}
