package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/domain"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/repo/dao"
	"github.com/samber/lo"
)

type FavoriteRepo interface {
	Create(ctx context.Context, favorite domain.Favorite) (int64, error)
	Update(ctx context.Context, favorite domain.Favorite) error
	Delete(ctx context.Context, id, uid int64) error
	FindByUid(ctx context.Context, biz string, uid int64) ([]domain.Favorite, error)
	CountUserFavorites(ctx context.Context, uid int64) (int64, error)
}

type NoCacheFavoriteRepo struct {
	dao dao.FavoriteDAO
}

func NewNoCacheFavoriteRepo(dao dao.FavoriteDAO) FavoriteRepo {
	return &NoCacheFavoriteRepo{
		dao: dao,
	}
}

func (repo *NoCacheFavoriteRepo) Create(ctx context.Context, favorite domain.Favorite) (int64, error) {
	return repo.dao.Insert(ctx, repo.domainToEntity(favorite))
}

func (repo *NoCacheFavoriteRepo) Update(ctx context.Context, favorite domain.Favorite) error {
	return repo.dao.Update(ctx, repo.domainToEntity(favorite))
}

func (repo *NoCacheFavoriteRepo) Delete(ctx context.Context, id, uid int64) error {
	return repo.dao.Delete(ctx, id, uid)
}

func (repo *NoCacheFavoriteRepo) FindByUid(ctx context.Context, biz string, uid int64) ([]domain.Favorite, error) {
	favorites, err := repo.dao.FindByUserId(ctx, biz, uid)
	if err != nil {
		return nil, err
	}
	return lo.Map(favorites, func(item dao.Favorite, index int) domain.Favorite {
		return repo.entityToDomain(item)
	}), nil
}

func (repo *NoCacheFavoriteRepo) domainToEntity(favorite domain.Favorite) dao.Favorite {
	return dao.Favorite{
		Id:        favorite.Id,
		UserId:    favorite.UserId,
		Name:      favorite.Name,
		Biz:       favorite.Biz,
		Private:   favorite.Private,
		CreatedAt: favorite.CreatedAt,
		UpdatedAt: favorite.UpdatedAt,
	}
}

func (repo *NoCacheFavoriteRepo) entityToDomain(favorite dao.Favorite) domain.Favorite {
	return domain.Favorite{
		Id:        favorite.Id,
		UserId:    favorite.UserId,
		Name:      favorite.Name,
		Biz:       favorite.Biz,
		Private:   favorite.Private,
		CreatedAt: favorite.CreatedAt,
		UpdatedAt: favorite.UpdatedAt,
	}
}
func (repo *NoCacheFavoriteRepo) CountUserFavorites(ctx context.Context, uid int64) (int64, error) {
	return repo.dao.CountUserFavorites(ctx, uid)
}
