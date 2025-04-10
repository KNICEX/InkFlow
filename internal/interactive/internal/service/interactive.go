package service

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/domain"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/events"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/repo"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"golang.org/x/sync/errgroup"
	"time"
)

type InteractiveService interface {
	CreateInteractive(ctx context.Context, biz string, bizId int64) error

	View(ctx context.Context, biz string, bizId int64, uid int64) error
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	Favorite(ctx context.Context, biz string, bizId int64, uid int64, fid int64) error

	CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error
	CancelFavorite(ctx context.Context, biz string, bizId int64, uid int64) error
	ListView(ctx context.Context, biz string, uid int64, maxId int64, limit int) ([]domain.ViewRecord, error)
	ListLike(ctx context.Context, biz string, uid int64, maxId int64, limit int) ([]domain.LikeRecord, error)
	ListFavoriteByFid(ctx context.Context, biz string, uid, fid int64, maxId int64, limit int) ([]domain.FavoriteRecord, error)

	CreateFavorite(ctx context.Context, favorite domain.Favorite) (int64, error)
	DeleteFavorite(ctx context.Context, favorite domain.Favorite) error
	FavoriteList(ctx context.Context, biz string, uid int64) ([]domain.Favorite, error)
	UpdateFavorite(ctx context.Context, favorite domain.Favorite) error

	Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error)
	GetMulti(ctx context.Context, biz string, bizIds []int64, uid int64) (map[int64]domain.Interactive, error)
	GetUserStats(ctx context.Context, uid int64) (domain.UserStats, error)
}

type interactiveService struct {
	repo     repo.InteractiveRepo
	favRepo  repo.FavoriteRepo
	producer events.InteractiveProducer
	l        logx.Logger
}

func NewInteractiveService(repo repo.InteractiveRepo, favRepo repo.FavoriteRepo, producer events.InteractiveProducer, l logx.Logger) InteractiveService {
	return &interactiveService{
		repo:     repo,
		producer: producer,
		favRepo:  favRepo,
		l:        l,
	}
}

// CreateInteractive
// TODO 在ink发布后需要创建
func (svc *interactiveService) CreateInteractive(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.CreateInteractive(ctx, biz, bizId)
}

func (svc *interactiveService) View(ctx context.Context, biz string, bizId int64, uid int64) error {
	if biz == domain.BizInk {
		go func() {
			if err := svc.producer.ProduceInkView(ctx, events.InkViewEvent{
				InkId:     bizId,
				UserId:    uid,
				CreatedAt: time.Now(),
			}); err != nil {
				svc.l.WithCtx(ctx).Error("produce ink view event error", logx.Error(err),
					logx.String("biz", biz),
					logx.Int64("bizId", bizId),
					logx.Int64("uid", uid))
			}
		}()
		return nil
	}
	return svc.repo.IncrView(ctx, biz, bizId, uid)
}

func (svc *interactiveService) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	if err := svc.repo.IncrLike(ctx, biz, bizId, uid); err != nil {
		return err
	}

	if biz == domain.BizInk {
		go func() {
			if err := svc.producer.ProduceInkLike(ctx, events.InkLikeEvent{
				InkId:     bizId,
				UserId:    uid,
				CreatedAt: time.Now(),
			}); err != nil {
				svc.l.WithCtx(ctx).Error("produce ink like event error", logx.Error(err),
					logx.String("biz", biz),
					logx.Int64("bizId", bizId),
					logx.Int64("uid", uid))
			}
		}()
	}
	return nil
}

func (svc *interactiveService) Favorite(ctx context.Context, biz string, bizId int64, uid int64, fid int64) error {
	return svc.repo.IncrFavorite(ctx, biz, bizId, uid, fid)

}

func (svc *interactiveService) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	if err := svc.repo.DecrLike(ctx, biz, bizId, uid); err != nil {
		return err
	}

	if biz == domain.BizInk {
		go func() {
			if err := svc.producer.ProduceInkCancelLike(ctx, events.InkCancelLikeEvent{
				InkId:     bizId,
				UserId:    uid,
				CreatedAt: time.Now(),
			}); err != nil {
				svc.l.WithCtx(ctx).Error("produce ink cancel like event error", logx.Error(err),
					logx.String("biz", biz),
					logx.Int64("bizId", bizId),
					logx.Int64("uid", uid))
			}
		}()
	}
	return nil
}

func (svc *interactiveService) CancelFavorite(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.DecrFavorite(ctx, biz, bizId, uid)
}

func (svc *interactiveService) ListView(ctx context.Context, biz string, uid int64, maxId int64, limit int) ([]domain.ViewRecord, error) {
	return svc.repo.ListView(ctx, biz, uid, maxId, limit)
}

func (svc *interactiveService) ListLike(ctx context.Context, biz string, uid int64, maxId int64, limit int) ([]domain.LikeRecord, error) {
	return svc.repo.ListLike(ctx, biz, uid, maxId, limit)
}

func (svc *interactiveService) ListFavoriteByFid(ctx context.Context, biz string, uid, fid int64, maxId int64, limit int) ([]domain.FavoriteRecord, error) {
	return svc.repo.ListFavoriteByFid(ctx, biz, uid, fid, maxId, limit)
}

func (svc *interactiveService) CreateFavorite(ctx context.Context, favorite domain.Favorite) (int64, error) {
	return svc.favRepo.Create(ctx, favorite)
}

func (svc *interactiveService) DeleteFavorite(ctx context.Context, favorite domain.Favorite) error {
	return svc.favRepo.Delete(ctx, favorite.Id, favorite.UserId)
}

func (svc *interactiveService) FavoriteList(ctx context.Context, biz string, uid int64) ([]domain.Favorite, error) {
	return svc.favRepo.FindByUid(ctx, biz, uid)
}

func (svc *interactiveService) UpdateFavorite(ctx context.Context, favorite domain.Favorite) error {
	return svc.favRepo.Update(ctx, favorite)
}

func (svc *interactiveService) Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error) {
	eg := errgroup.Group{}
	var intr domain.Interactive
	eg.Go(func() error {
		var er error
		intr, er = svc.repo.Get(ctx, biz, bizId)
		return er
	})
	var liked bool
	var favorited bool
	if uid != 0 {
		// 聚合是否点赞
		// TODO 后续还有收藏或更多其他操作
		eg.Go(func() error {
			var er error
			liked, er = svc.repo.Liked(ctx, biz, bizId, uid)
			return er
		})
		eg.Go(func() error {
			var er error
			favorited, er = svc.repo.Favorited(ctx, biz, bizId, uid)
			return er
		})
	}
	if err := eg.Wait(); err != nil {
		return domain.Interactive{}, err
	}
	intr.Liked = liked
	intr.Favorited = favorited
	return intr, nil
}

func (svc *interactiveService) GetMulti(ctx context.Context, biz string, bizIds []int64, uid int64) (map[int64]domain.Interactive, error) {
	var likedMap map[int64]bool
	var favoritedMap map[int64]bool

	// 加载点赞和收藏状态
	eg := errgroup.Group{}
	if uid != 0 {
		eg.Go(func() error {
			var er error
			likedMap, er = svc.repo.LikedBatch(ctx, biz, bizIds, uid)
			return er
		})
		eg.Go(func() error {
			var er error
			favoritedMap, er = svc.repo.FavoritedBatch(ctx, biz, bizIds, uid)
			return er
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	intrs, err := svc.repo.GetBatch(ctx, biz, bizIds)
	if err != nil {
		return nil, err
	}

	for _, intr := range intrs {
		if liked, ok := likedMap[intr.BizId]; ok {
			intr.Liked = liked
		}
		if favorited, ok := favoritedMap[intr.BizId]; ok {
			intr.Favorited = favorited
		}
	}
	return intrs, nil
}

func (svc *interactiveService) GetUserStats(ctx context.Context, uid int64) (domain.UserStats, error) {
	var favCount, likeCount, viewCount int64
	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		favCount, er = svc.favRepo.CountUserFavorites(ctx, uid)
		return er
	})
	eg.Go(func() error {
		var er error
		likeCount, er = svc.repo.CountUserLikes(ctx, uid)
		return er
	})
	eg.Go(func() error {
		var er error
		viewCount, er = svc.repo.CountUserViews(ctx, uid)
		return er
	})
	if err := eg.Wait(); err != nil {
		return domain.UserStats{}, err
	}

	return domain.UserStats{
		FavoriteCnt: favCount,
		LikeCnt:     likeCount,
		ViewCnt:     viewCount,
	}, nil
}
