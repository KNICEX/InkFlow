//go:build wireinject

package review

import (
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/ai"
	"github.com/KNICEX/InkFlow/internal/review/internal/event"
	"github.com/KNICEX/InkFlow/internal/review/internal/repo"
	"github.com/KNICEX/InkFlow/internal/review/internal/repo/dao"
	"github.com/KNICEX/InkFlow/internal/review/internal/service/failover"
	"github.com/KNICEX/InkFlow/internal/review/internal/service/llm"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"github.com/google/wire"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
)

func InitAsyncService(producer sarama.SyncProducer, l logx.Logger) AsyncService {
	wire.Build(
		event.NewKafkaReviewProducer,
		llm.NewAsyncWorkflowService,
	)
	return nil
}

func InitService(llmSvc ai.LLMService) Service {
	return llm.NewLLMService(llmSvc)
}

func InitReviewConsumer(workflowCli client.Client, saramaCli sarama.Client, service Service, failoverSvc FailoverService, l logx.Logger) *event.ReviewConsumer {
	wire.Build(
		event.NewReviewConsumer,
	)
	return nil
}

func initSnowflakeNode() snowflakex.Node {
	return snowflakex.NewNode(snowflakex.DefaultStartTime, 0)
}

func initFailoverDao(db *gorm.DB) dao.ReviewFailDAO {
	if err := dao.InitTables(db); err != nil {
		panic(err)
	}

	return dao.NewGormReviewFailDAO(db, initSnowflakeNode())
}

func InitFailoverService(workflowCli client.Client, svc Service, db *gorm.DB, l logx.Logger) FailoverService {
	wire.Build(
		initFailoverDao,
		repo.NewReviewFailRepo,
		failover.NewReviewService,
	)
	return nil
}
