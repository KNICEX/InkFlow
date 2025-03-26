//go:build wireinject

package search

import (
	"github.com/KNICEX/InkFlow/internal/search/internal/repo"
	"github.com/KNICEX/InkFlow/internal/search/internal/repo/dao"
	"github.com/KNICEX/InkFlow/internal/search/internal/service"
	"github.com/google/wire"
	"github.com/meilisearch/meilisearch-go"
)

func InitSearchService() Service {
	return nil
}

func initInkDAO(meili meilisearch.ServiceManager) dao.InkDAO {
	if err := dao.InitMeili(meili); err != nil {
		panic(err)
	}
	return dao.NewMeiliInkDAO(meili)
}

func InitSyncService(meili meilisearch.ServiceManager) SyncService {
	wire.Build(
		dao.NewMeiliCommentDAO,
		dao.NewMeiliUserDAO,
		initInkDAO,
		repo.NewCommentRepo,
		repo.NewUserRepo,
		repo.NewInkRepo,
		service.NewSyncService)
	return nil
}
