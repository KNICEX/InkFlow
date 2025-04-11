package search

import (
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/search/internal/event"
	"github.com/KNICEX/InkFlow/internal/search/internal/repo"
	"github.com/KNICEX/InkFlow/internal/search/internal/repo/dao"
	"github.com/KNICEX/InkFlow/internal/search/internal/service"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/meilisearch/meilisearch-go"
	"sync"
)

var (
	inkRepo     repo.InkRepo
	commentRepo repo.CommentRepo
	userRepo    repo.UserRepo
	once        sync.Once
)

func initRepo(meili meilisearch.ServiceManager) {
	once.Do(func() {
		if err := dao.InitMeili(meili); err != nil {
			panic(err)
		}
		userDao := dao.NewMeiliUserDAO(meili)
		inkRepo = repo.NewInkRepo(dao.NewMeiliInkDAO(meili), userDao)
		commentRepo = repo.NewCommentRepo(dao.NewMeiliCommentDAO(meili), userDao)
		userRepo = repo.NewUserRepo(userDao)
	})
}

func InitSearchService(meili meilisearch.ServiceManager) Service {
	initRepo(meili)
	return service.NewSearchService(userRepo, inkRepo, commentRepo)
}

func InitSyncConsumer(svc SyncService, cli sarama.Client, l logx.Logger) *SyncConsumer {
	replyHandler := event.NewReplyHandler(svc)
	userCreateHandler := event.NewUserCreateHandler(svc)
	userUpdateHandler := event.NewUserUpdateHandler(svc)

	consumer := event.NewSyncConsumer(cli, svc, l)
	if err := consumer.RegisterHandler(replyHandler, userCreateHandler, userUpdateHandler); err != nil {
		panic(err)
	}
	return consumer
}

func InitSyncService(meili meilisearch.ServiceManager) SyncService {
	initRepo(meili)
	return service.NewSyncService(userRepo, inkRepo, commentRepo)
}
