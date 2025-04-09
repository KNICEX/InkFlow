package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/KNICEX/InkFlow/internal/user/internal/domain"
	"github.com/KNICEX/InkFlow/internal/user/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/user/internal/repo/dao"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/samber/lo"
	"time"
)

var (
	ErrUserNotFound  = dao.ErrRecordNotFound
	ErrUserDuplicate = dao.ErrDuplicateKey
)

type UserRepo interface {
	Create(ctx context.Context, u domain.User) (int64, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByAccount(ctx context.Context, accountName string) (domain.User, error)

	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByIds(ctx context.Context, ids []int64) (map[int64]domain.User, error)
	FindByGithubId(ctx context.Context, id int64) (domain.User, error)
	UpdateNonZeroFields(ctx context.Context, u domain.User) error
}

var _ UserRepo = (*CachedUserRepo)(nil)

type CachedUserRepo struct {
	dao   dao.UserDAO
	cache cache.UserCache
	l     logx.Logger
}

func NewCachedUserRepo(dao dao.UserDAO, cache cache.UserCache, l logx.Logger) UserRepo {
	return &CachedUserRepo{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (r *CachedUserRepo) Create(ctx context.Context, u domain.User) (int64, error) {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *CachedUserRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepo) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepo) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(context.Background(), id)
	if err == nil {
		return u, nil
	}
	//if errors.Is(err, cache.ErrKeyNotExist) {
	//	// 缓存未命中
	//}
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = r.entityToDomain(ue)
	go func() {
		err = r.cache.Set(context.WithoutCancel(ctx), id, u)
		if err != nil {
			r.l.Error("set user cache error", logx.Error(err), logx.Int64("UserId", id))
		}
	}()
	return u, err
}

func (r *CachedUserRepo) FindByIds(ctx context.Context, ids []int64) (map[int64]domain.User, error) {
	users, err := r.cache.GetByIds(ctx, ids)
	if err != nil && !errors.Is(err, cache.ErrKeyNotExist) {
		r.l.WithCtx(ctx).Error("get user cache by ids error", logx.Error(err), logx.Any("UserIds", ids))
	}

	if len(users) == len(ids) {
		return users, nil
	}

	if len(users) > 0 {
		// 去除缓存命中的
		ids = lo.Reject(ids, func(item int64, index int) bool {
			_, ok := users[item]
			return ok
		})
	}

	// 从数据库中获取
	us, err := r.dao.FindByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, u := range us {
		users[u.Id] = r.entityToDomain(u)
	}

	return users, nil
}

func (r *CachedUserRepo) FindByWechatOpenId(ctx context.Context, openId string) (domain.User, error) {
	u, err := r.dao.FindByWechatOpenId(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepo) FindByGithubId(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.dao.FindByGithubId(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepo) UpdateNonZeroFields(ctx context.Context, u domain.User) error {
	if err := r.dao.UpdateById(ctx, r.domainToEntity(u)); err != nil {
		return err
	}
	// 延时双删
	defer time.AfterFunc(time.Second*3, func() {
		er := r.cache.Delete(ctx, u.Id)
		if er != nil {
			r.l.Error("delayed delete user cache error", logx.Error(er), logx.Int64("UserId", u.Id))
		}
	})

	if err := r.cache.Delete(ctx, u.Id); err != nil {
		// 删除缓存错误，不认为操作失败
		r.l.Error("delete user cache error", logx.Error(err), logx.Int64("UserId", u.Id))
	}
	return nil
}

func (r *CachedUserRepo) FindByAccount(ctx context.Context, accountName string) (domain.User, error) {
	u, err := r.dao.FindByAccountName(ctx, accountName)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepo) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:         u.Id,
		Email:      u.Email.String,
		Phone:      u.Phone.String,
		Password:   u.Password.String,
		Username:   u.Username,
		Account:    u.Account,
		AboutMe:    u.AboutMe,
		Avatar:     u.Avatar,
		Banner:     u.Banner,
		Links:      domain.LinksFromString(u.Links),
		GithubInfo: u.GithubInfo,
		Birthday:   u.Birthday.Time,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func (r *CachedUserRepo) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},

		Account:  u.Account,
		Avatar:   u.Avatar,
		Banner:   u.Banner,
		Username: u.Username,
		AboutMe:  u.AboutMe,
		Links:    u.Links.ToString(),
		Password: sql.NullString{
			String: u.Password,
			Valid:  u.Password != "",
		},
		GithubId: sql.NullInt64{
			Int64: u.GithubInfo.Id,
			Valid: u.GithubInfo.Id != 0,
		},
		GithubInfo: u.GithubInfo,

		Birthday: sql.NullTime{
			Time:  u.Birthday,
			Valid: !u.Birthday.IsZero(),
		},
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
