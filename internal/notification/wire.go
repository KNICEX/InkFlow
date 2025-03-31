//go:build wireinject

package notification

import (
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/comment"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/notification/internal/event"
	"github.com/KNICEX/InkFlow/internal/notification/internal/repo"
	"github.com/KNICEX/InkFlow/internal/notification/internal/repo/dao"
	"github.com/KNICEX/InkFlow/internal/notification/internal/service"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func initSnowflakeNode() snowflakex.Node {
	return snowflakex.NewNode(snowflakex.DefaultStartTime, 0)
}

func initDAO(db *gorm.DB, node snowflakex.Node) dao.NotificationDAO {
	if err := dao.InitTables(db); err != nil {
		panic(err)
	}
	return dao.NewGormNotificationDAO(db, node)
}

func InitNotificationService(db *gorm.DB) Service {
	wire.Build(
		initSnowflakeNode,
		initDAO,
		repo.NewNoCacheNotificationRepo,
		service.NewNotificationService,
	)
	return nil
}

func InitNotificationConsumer(cli sarama.Client, svc Service, inkSvc ink.Service, commentSvc comment.Service, l logx.Logger) *SyncConsumer {
	consumer := event.NewNotificationConsumer(cli, svc, l)
	userFollowHandler := event.NewFollowHandler(svc)
	commentReplyHandler := event.NewReplyHandler(svc, commentSvc, inkSvc)
	commentLikeHandler := event.NewCommentLikeHandler(svc, commentSvc)
	inkLikeHandler := event.NewInkLikeHandler(svc, inkSvc)

	err := consumer.RegisterHandler(userFollowHandler, commentReplyHandler, commentLikeHandler, inkLikeHandler)
	if err != nil {
		panic(err)
	}
	return consumer
}
