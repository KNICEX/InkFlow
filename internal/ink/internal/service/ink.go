package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/internal/ink/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/yuin/goldmark"
)

var (
	ErrNotFound     = errors.New("ink not found")
	ErrNoPermission = errors.New("no permission")
)

// InkService
// 这里将Save和Publish分开，Sava只当作保存草稿，前端编辑时也会定时调用放丢失
// 直接发布文章，前端先调用Save保存草稿，然后再调用Publish发布
type InkService interface {
	Save(ctx context.Context, ink domain.Ink) (int64, error) // 保存草稿
	Publish(ctx context.Context, ink domain.Ink) error       // 发布

	UpdateLiveStatus(ctx context.Context, id int64, authorId int64, status domain.Status) error // 更新文章状态
	UpdateDraftStatus(ctx context.Context, id int64, authorId int64, status domain.Status) error
	SyncToLive(ctx context.Context, ink domain.Ink) error // 同步到线上库

	Withdraw(ctx context.Context, ink domain.Ink) error    // 撤回
	DeleteDraft(ctx context.Context, ink domain.Ink) error // 删除草稿
	DeleteLive(ctx context.Context, ink domain.Ink) error  // 删除线上文章

	FindByIds(ctx context.Context, ids []int64) (map[int64]domain.Ink, error) // 批量获取文章
	FindLiveInk(ctx context.Context, id int64) (domain.Ink, error)            // 获取公开ink
	FindDraftInk(ctx context.Context, id, authorId int64) (domain.Ink, error) // 获取草稿ink
	FindPendingInk(ctx context.Context, id, authorId int64) (domain.Ink, error)
	FindPrivateInk(ctx context.Context, id, authorId int64) (domain.Ink, error)  // 获取私有ink
	FindRejectedInk(ctx context.Context, id, authorId int64) (domain.Ink, error) // 获取审核拒绝的ink
	ListLiveByAuthorId(ctx context.Context, authorId int64, offset int, limit int) ([]domain.Ink, error)
	ListPendingByAuthorId(ctx context.Context, authorId int64, offset int, limit int) ([]domain.Ink, error)
	ListReviewRejectedByAuthorId(ctx context.Context, authorId int64, offset int, limit int) ([]domain.Ink, error)
	ListDraftByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error)
	ListPrivateByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error)
	ListAllLive(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error)
	ListAllReviewRejected(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error)
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

	// 将markdown转换为html
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(ink.ContentMeta), &buf); err != nil {
		return 0, err
	}
	ink.ContentHtml = buf.String()

	if ink.Id == 0 {
		return svc.draftRepo.Create(ctx, ink)
	}
	return ink.Id, svc.draftRepo.Update(ctx, ink)
}

// Publish
// 从草稿库查询文章(因为更新文章都是先save更新草稿，我们认为草稿是最新的数据)
func (svc *inkService) Publish(ctx context.Context, ink domain.Ink) error {
	draft, err := svc.draftRepo.FindByIdAndAuthorId(ctx, ink.Id, ink.Author.Id)
	if err != nil {
		return err
	}
	if draft.Status != domain.InkStatusUnPublished {
		svc.l.WithCtx(ctx).Error("InkService Publish  draft status error",
			logx.Error(err), logx.Int64("inkId", ink.Id), logx.Any("status", draft.Status))
		return errors.New("invalid status")
	}

	// 修改草稿的状态为待审核
	return svc.draftRepo.UpdateStatus(ctx, domain.Ink{
		Id: draft.Id,
		Author: domain.Author{
			Id: draft.Author.Id,
		},
		Status: domain.InkStatusPending,
	})
}

// UpdateInkStatus
// 此方法使用权不应该交给用户，仅内部服务可以调用
func (svc *inkService) UpdateLiveStatus(ctx context.Context, id int64, authorId int64, status domain.Status) error {
	return svc.liveRepo.UpdateStatus(ctx, domain.Ink{
		Id: id,
		Author: domain.Author{
			Id: authorId,
		},
		Status: status,
	})
}

func (svc *inkService) UpdateDraftStatus(ctx context.Context, id int64, authorId int64, status domain.Status) error {
	return svc.draftRepo.UpdateStatus(ctx, domain.Ink{
		Id: id,
		Author: domain.Author{
			Id: authorId,
		},
		Status: status,
	})
}

func (svc *inkService) SyncToLive(ctx context.Context, ink domain.Ink) error {
	_, err := svc.liveRepo.Save(ctx, ink)
	return err
}

func (svc *inkService) Withdraw(ctx context.Context, ink domain.Ink) error {
	// TODO 暂不可用
	ink.Status = domain.InkStatusPrivate
	err := svc.liveRepo.UpdateStatus(ctx, ink)
	if err != nil {
		return err
	}
	// 更新草稿状态, 让草稿可以继续编辑
	err = svc.draftRepo.UpdateStatus(ctx, domain.Ink{
		Id: ink.Id,
		Author: domain.Author{
			Id: ink.Author.Id,
		},
		Status: domain.InkStatusPublished,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *inkService) DeleteDraft(ctx context.Context, ink domain.Ink) error {
	// 针对草稿，只能删除未发布的草稿，同时把线上库的文章(此时线上库的文章的状态应该为Unpublished)也删除
	err := svc.draftRepo.Delete(ctx, ink.Id, ink.Author.Id, domain.InkStatusUnPublished)
	if err != nil {
		return err
	}
	return svc.liveRepo.Delete(ctx, ink.Id, ink.Author.Id, domain.InkStatusUnPublished)
}

func (svc *inkService) DeleteLive(ctx context.Context, ink domain.Ink) error {
	err := svc.liveRepo.Delete(ctx, ink.Id, ink.Author.Id, domain.InkStatusPublished, domain.InkStatusRejected, domain.InkStatusPrivate)
	if err != nil {
		return err
	}
	// 删除线上，草稿直接一起删除
	err = svc.draftRepo.Delete(ctx, ink.Id, ink.Author.Id)
	return err
}

func (svc *inkService) wrapNotFoundErr(err error) error {
	if errors.Is(err, repo.ErrLiveInkNotFound) || errors.Is(err, repo.ErrDraftNotFound) {
		return ErrNotFound
	}
	return err
}

func (svc *inkService) FindByIds(ctx context.Context, ids []int64) (map[int64]domain.Ink, error) {
	inks, err := svc.liveRepo.FindByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	inkMap := make(map[int64]domain.Ink)
	for _, ink := range inks {
		inkMap[ink.Id] = ink
	}
	return inkMap, nil
}

func (svc *inkService) FindLiveInk(ctx context.Context, id int64) (domain.Ink, error) {
	ink, err := svc.liveRepo.FindById(ctx, id, domain.InkStatusPublished)
	if err != nil {
		return domain.Ink{}, svc.wrapNotFoundErr(err)
	}
	return ink, nil
}

func (svc *inkService) FindDraftInk(ctx context.Context, id, authorId int64) (domain.Ink, error) {
	draft, err := svc.draftRepo.FindByIdAndAuthorId(ctx, id, authorId, domain.InkStatusUnPublished)
	if err != nil {
		return domain.Ink{}, svc.wrapNotFoundErr(err)
	}
	return draft, nil
}
func (svc *inkService) FindPendingInk(ctx context.Context, id, authorId int64) (domain.Ink, error) {
	ink, err := svc.draftRepo.FindByIdAndAuthorId(ctx, id, authorId, domain.InkStatusPending)
	if err != nil {
		return domain.Ink{}, svc.wrapNotFoundErr(err)
	}
	return ink, nil
}
func (svc *inkService) FindRejectedInk(ctx context.Context, id, authorId int64) (domain.Ink, error) {
	//return svc.draftRepo.FindByIdAndAuthorId(ctx, id, authorId, domain.InkStatusRejected)
	ink, err := svc.liveRepo.FindById(ctx, id, domain.InkStatusRejected)
	if err != nil {
		return domain.Ink{}, svc.wrapNotFoundErr(err)
	}
	return ink, nil
}

func (svc *inkService) FindPrivateInk(ctx context.Context, id, authorId int64) (domain.Ink, error) {
	ink, err := svc.liveRepo.FindById(ctx, id, domain.InkStatusPrivate)
	if err != nil {
		return domain.Ink{}, svc.wrapNotFoundErr(err)
	}
	if ink.Author.Id != authorId {
		return domain.Ink{}, ErrNoPermission
	}
	return ink, nil
}

func (svc *inkService) ListLiveByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error) {
	return svc.liveRepo.FindByAuthorId(ctx, authorId, offset, limit, domain.InkStatusPublished)
}

func (svc *inkService) ListPendingByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error) {
	return svc.draftRepo.FindByAuthorId(ctx, authorId, offset, limit, domain.InkStatusPending)
}

func (svc *inkService) ListReviewRejectedByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error) {
	return svc.draftRepo.FindByAuthorId(ctx, authorId, offset, limit, domain.InkStatusRejected)
}

func (svc *inkService) ListDraftByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error) {
	return svc.draftRepo.FindByAuthorId(ctx, authorId, offset, limit, domain.InkStatusUnPublished)
}

func (svc *inkService) ListPrivateByAuthorId(ctx context.Context, authorId int64, offset, limit int) ([]domain.Ink, error) {
	return svc.liveRepo.FindByAuthorId(ctx, authorId, offset, limit, domain.InkStatusPrivate)
}

// ListAllLive 不要暴露给用户
func (svc *inkService) ListAllLive(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error) {
	return svc.liveRepo.FindAll(ctx, maxId, limit, domain.InkStatusPublished)
}

func (svc *inkService) ListAllReviewRejected(ctx context.Context, maxId int64, limit int) ([]domain.Ink, error) {
	return svc.liveRepo.FindAll(ctx, maxId, limit, domain.InkStatusRejected)
}

type InkServiceV1 interface {
	FindByAuthorAndStatus(ctx context.Context, authorId int64, status domain.Status, offset, limit int) ([]domain.Ink, error)
}
