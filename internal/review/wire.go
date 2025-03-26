//go:build wireinject

package review

import (
	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/internal/ai"
	"github.com/KNICEX/InkFlow/internal/review/internal/event"
	"github.com/KNICEX/InkFlow/internal/review/internal/service/llm"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"github.com/google/wire"
	"go.temporal.io/sdk/client"
)

func InitAsyncService(producer sarama.SyncProducer, l logx.Logger) AsyncService {
	wire.Build(
		event.NewKafkaReviewProducer,
		llm.NewAsyncWorkflowService,
	)
	return nil
}

func InitReviewConsumer(workflowCli client.Client, saramaCli sarama.Client, service ai.LLMService, l logx.Logger) *event.ReviewConsumer {
	wire.Build(
		llm.NewLLMService,
		event.NewReviewConsumer,
	)
	return nil
}
