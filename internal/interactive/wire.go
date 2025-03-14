//go:build wireinject

package interactive

import (
	"github.com/KNICEX/InkFlow/internal/interactive/internal/repo"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/repo/dao"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/service"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func initSnowflakeNode() snowflakex.Node {
	return snowflakex.NewNode(snowflakex.DefaultStartTime, 0)
}

func initDAO(db *gorm.DB, node snowflakex.Node, l logx.Logger) dao.InteractiveDAO {
	if err := dao.InitTables(db); err != nil {
		panic(err)
	}
	return dao.NewGormInteractiveDAO(db, node, l)
}

func InitInteractiveService(cmd redis.Cmdable, db *gorm.DB, l logx.Logger) Service {
	wire.Build(
		initSnowflakeNode,
		dao.NewGormInteractiveDAO,
		cache.NewRedisInteractiveCache,
		repo.NewCachedInteractiveRepo,
		service.NewInteractiveService,
	)
	return nil
}
