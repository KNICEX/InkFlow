//go:build wireinject

package ink

import (
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/ink/internal/repo/dao"
	"github.com/KNICEX/InkFlow/internal/ink/internal/service"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func initSnowflakeNode() snowflakex.Node {
	return snowflakex.NewNode(snowflakex.DefaultStartTime, 0)
}

func initDraftDAO(db *gorm.DB, node snowflakex.Node) dao.DraftDAO {
	if err := dao.InitTable(db); err != nil {
		panic(err)
	}
	return dao.NewDraftDAO(db, node)
}

func InitInkService(cmd redis.Cmdable, db *gorm.DB, l logx.Logger) Service {
	wire.Build(
		initSnowflakeNode,
		initDraftDAO,
		dao.NewLiveDAO,
		cache.NewRedisInkCache,
		repo.NewCachedLiveInkRepo,
		repo.NewNoCacheDraftInkRepo,
		service.NewInkService,
	)
	return nil
}
