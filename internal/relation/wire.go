//go:build wireinject

package relation

import (
	"github.com/KNICEX/InkFlow/internal/relation/internal/repo"
	"github.com/KNICEX/InkFlow/internal/relation/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/relation/internal/repo/dao"
	"github.com/KNICEX/InkFlow/internal/relation/internal/service"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func initSnowflake() snowflakex.Node {
	return snowflakex.NewNode(snowflakex.DefaultStartTime, 0)
}

func initFollowDAO(db *gorm.DB, node snowflakex.Node, l logx.Logger) dao.FollowRelationDAO {
	if err := dao.InitTable(db); err != nil {
		panic(err)
	}
	return dao.NewGormFollowRelationDAO(db, node, l)
}

func InitFollowService(cmd redis.Cmdable, db *gorm.DB, l logx.Logger) FollowService {
	wire.Build(
		initSnowflake,
		initFollowDAO,
		cache.NewRedisFollowCache,
		repo.NewCachedFollowRepo,
		service.NewFollowService,
	)
	return nil
}
