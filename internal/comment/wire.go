//go:build wireinject

package comment

import (
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/comment/internal/event"
	"github.com/KNICEX/InkFlow/internal/comment/internal/repo"
	"github.com/KNICEX/InkFlow/internal/comment/internal/repo/cache"
	"github.com/KNICEX/InkFlow/internal/comment/internal/repo/dao"
	"github.com/KNICEX/InkFlow/internal/comment/internal/service"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func initSnowflakeNode() snowflakex.Node {
	return snowflakex.NewNode(snowflakex.DefaultStartTime, 0)
}

func initDAO(db *gorm.DB, node snowflakex.Node, l logx.Logger) dao.CommentDAO {
	dao.Init(db)
	return dao.NewGormCommentDAO(db, node, l)
}

func InitCommentService(db *gorm.DB, cmd redis.Cmdable, inkSvc ink.Service, producer sarama.SyncProducer, l logx.Logger) Service {
	wire.Build(
		initSnowflakeNode,
		initDAO,
		cache.NewRedisCommentCache,
		repo.NewCachedCommentRepo,
		event.NewKafkaCommentEvtProducer,
		service.NewCommentService)
	return nil
}
