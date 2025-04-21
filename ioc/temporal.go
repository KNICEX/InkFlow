package ioc

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/workflow/inkpub"
	"github.com/KNICEX/InkFlow/internal/workflow/schedule"
	"github.com/KNICEX/InkFlow/pkg/schedulex"
	"github.com/KNICEX/InkFlow/pkg/temporalx"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	inkPubQueue      = "ink-pub-queue"
	rankInkQueue     = "rank-ink-queue"
	rankTagQueue     = "rank-tag-queue"
	retryReviewQueue = "retry-review-queue"
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

type RankInkWorker struct {
	worker.Worker
}

type RankTagWorker struct {
	worker.Worker
}

type RetryReviewWorker struct {
	worker.Worker
}

func InitRankInkWorker(cli client.Client, activities *schedule.RankActivities) *RankInkWorker {
	w := worker.New(cli, rankInkQueue, worker.Options{})
	w.RegisterWorkflow(schedule.RankHotInk)
	w.RegisterActivity(activities)
	return &RankInkWorker{
		Worker: w,
	}
}

func InitRankTagWorker(cli client.Client, activities *schedule.RankActivities) *RankTagWorker {
	w := worker.New(cli, rankTagQueue, worker.Options{})
	w.RegisterWorkflow(schedule.RankHotTag)
	w.RegisterActivity(activities)
	return &RankTagWorker{
		Worker: w,
	}
}

func InitRetryReviewWorker(cli client.Client, activities *schedule.ReviewFailoverActivity) *RetryReviewWorker {
	w := worker.New(cli, retryReviewQueue, worker.Options{})
	w.RegisterWorkflow(schedule.ReviewRetryFail)
	w.RegisterActivity(activities)
	return &RetryReviewWorker{
		Worker: w,
	}
}

func InitWorkers(inkPub *InkPubWorker, rankTag *RankTagWorker, rankInk *RankInkWorker, retryReview *RetryReviewWorker) []worker.Worker {
	return []worker.Worker{
		inkPub.Worker,
		rankTag.Worker,
		rankInk.Worker,
		retryReview.Worker,
	}
}

type RankInkScheduler func() error

func (r RankInkScheduler) Start() error {
	return r()
}

func InitRankInkScheduler(cli client.Client) RankInkScheduler {
	return func() error {
		return temporalx.UpsertSchedule(context.Background(), cli, client.ScheduleOptions{
			ID: "rank-ink-scheduler",
			Spec: client.ScheduleSpec{
				CronExpressions: []string{"@every 10m"},
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        "rank-ink-scheduler-action",
				Workflow:  schedule.RankHotInk,
				TaskQueue: rankInkQueue,
			},
		})
	}
}

type RankTagScheduler func() error

func (r RankTagScheduler) Start() error {
	return r()
}

func InitRankTagScheduler(cli client.Client) RankTagScheduler {
	return func() error {
		return temporalx.UpsertSchedule(context.Background(), cli, client.ScheduleOptions{
			ID: "rank-tag-scheduler",
			Spec: client.ScheduleSpec{
				CronExpressions: []string{"@every 30m"},
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        "rank-tag-scheduler-action",
				Workflow:  schedule.RankHotTag,
				TaskQueue: rankTagQueue,
			},
		})
	}
}

type ReviewFailRetryScheduler func() error

func (r ReviewFailRetryScheduler) Start() error {
	return r()
}

func InitReviewRetryScheduler(cli client.Client) ReviewFailRetryScheduler {
	return func() error {
		return temporalx.UpsertSchedule(context.Background(), cli, client.ScheduleOptions{
			ID: "review-fail-retry-scheduler",
			Spec: client.ScheduleSpec{
				CronExpressions: []string{"@every 5m"},
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        "review-fail-retry-scheduler-action",
				Workflow:  schedule.ReviewRetryFail,
				TaskQueue: retryReviewQueue,
			},
		})
	}
}

func InitSchedulers(rankInk RankInkScheduler, rankTag RankTagScheduler, reviewRetry ReviewFailRetryScheduler) []schedulex.Scheduler {
	return []schedulex.Scheduler{
		rankInk,
		rankTag,
		reviewRetry,
	}
}
