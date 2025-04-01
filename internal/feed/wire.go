package feed

import (
	"github.com/KNICEX/InkFlow/internal/action"
	"github.com/KNICEX/InkFlow/internal/feed/internal/repo"
	"github.com/KNICEX/InkFlow/internal/feed/internal/repo/dao"
	"github.com/KNICEX/InkFlow/internal/feed/internal/service"
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"gorm.io/gorm"
)

var r repo.FeedRepo

func initRepo(db *gorm.DB, l logx.Logger) repo.FeedRepo {
	if r != nil {
		return r
	}
	if err := dao.InitTables(db); err != nil {
		panic(err)
	}
	node := snowflakex.NewNode(snowflakex.DefaultStartTime, 0)
	pushDAO := dao.NewGormPushFeedDAO(db, node)
	pullDAO := dao.NewGormFeedPullDAO(db, node)
	r = repo.NewNoCacheFeedRepo(pushDAO, pullDAO, l)
	return r
}

func InitService(db *gorm.DB, followSvc relation.FollowService, actionSvc action.Service, l logx.Logger) Service {
	return service.NewFeedService(initRepo(db, l), followSvc, actionSvc)
}
