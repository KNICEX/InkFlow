//go:build wireinject

package interactive

import (
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/interactive/internal/events"
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

// 为了初始化consumer,不得已使用一个包变量实现单例
var r repo.InteractiveRepo

func initIntrRepo(d dao.InteractiveDAO, c cache.InteractiveCache, l logx.Logger) repo.InteractiveRepo {
	if r != nil {
		return r
	}
	r = repo.NewCachedInteractiveRepo(c, d, l)
	return r
}
func initProducer(p sarama.SyncProducer) events.InteractiveProducer {
	return events.NewKafkaInteractiveProducer(p)
}

func InitInteractiveService(cmd redis.Cmdable, p sarama.SyncProducer, db *gorm.DB, l logx.Logger) Service {
	wire.Build(
		initSnowflakeNode,
		dao.NewGormInteractiveDAO,
		cache.NewRedisInteractiveCache,
		initProducer,
		initIntrRepo,
		service.NewInteractiveService,
	)
	return nil
}

func InitInteractiveInkReadConsumer(client sarama.Client, l logx.Logger) *events.InkViewConsumer {
	return events.NewInkViewConsumer(client, r, l)
}
