package user

import (
	"github.com/KNICEX/InkFlow/internal/user/internal/repo"
	"github.com/KNICEX/InkFlow/internal/user/internal/repo/dao"
	"github.com/KNICEX/InkFlow/internal/user/internal/service"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitUserService(db *gorm.DB, cmd redis.Cmdable, l logx.Logger) Service {
	wire.Build(
		dao.NewGormUserDAO,
		repo.NewCachedUserRepo,
		service.NewUserService,
	)
	return service.NewUserService(nil, nil)
}
