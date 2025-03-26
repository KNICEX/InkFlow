package ioc

import (
	"github.com/KNICEX/InkFlow/internal/workflow/inkpub"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	inkPubQueue = "ink-pub-queue"
)

func InitTemporalClient() client.Client {
	type Config struct {
		Addr string `mapstructure:"addr"`
	}
	var cfg Config
	if err := viper.UnmarshalKey("temporal", &cfg); err != nil {
		panic(err)
	}
	cli, err := client.Dial(client.Options{
		HostPort: cfg.Addr,
	})
	if err != nil {
		panic(err)
	}
	return cli
}

type InkPubWorker struct {
	worker.Worker
}

func InitInkPubWorker(cli client.Client, activities *inkpub.Activities) *InkPubWorker {
	w := worker.New(cli, inkPubQueue, worker.Options{})
	w.RegisterWorkflow(inkpub.InkPublish)
	w.RegisterActivity(activities)
	return &InkPubWorker{
		Worker: w,
	}
}

func InitWorkers(inkPub *InkPubWorker) []worker.Worker {
	return []worker.Worker{
		inkPub.Worker,
	}
}
