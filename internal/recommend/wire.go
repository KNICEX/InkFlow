package recommend

import (
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/interactive"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/event"
	"github.com/KNICEX/InkFlow/internal/recommend/internal/service/gorse"
	"github.com/KNICEX/InkFlow/internal/relation"
	"github.com/KNICEX/InkFlow/pkg/logx"
	client "github.com/gorse-io/gorse-go"
)

func InitSyncService(cli *client.GorseClient) SyncService {
	return gorse.NewSyncService(cli)
}

func InitSyncConsumer(cli sarama.Client, svc SyncService, l logx.Logger) *SyncConsumer {
	userCreateHandler := event.NewUserCreateHandler(svc)
	consumer := event.NewSyncConsumer(cli, l)
	if err := consumer.RegisterHandler(userCreateHandler); err != nil {
		panic(err)
	}
	return consumer
}

func InitService(cli *client.GorseClient, followSvc relation.FollowService, intrSvc interactive.Service, l logx.Logger) Service {
	return gorse.NewRecommendService(cli, followSvc, intrSvc, l)
}
