package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/domain"
	"github.com/KNICEX/InkFlow/internal/repo/dao"
)

var (
	ErrUserNotFound  = dao.ErrRecordNotFound
	ErrUserDuplicate = dao.ErrDuplicateKey
)

type UserRepo interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByAccountName(ctx context.Context, accountName string) (domain.User, error)

	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByGithubId(ctx context.Context, id int64) (domain.User, error)
	UpdateNonZeroFields(ctx context.Context, u domain.User) error
}
