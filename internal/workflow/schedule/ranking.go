package schedule

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func RankHotInk(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: time.Second * 5,
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var activities *RankActivities
	l := workflow.GetLogger(ctx)
	err := workflow.ExecuteActivity(ctx, activities.RankInk, 1000).Get(ctx, nil)
	if err != nil {
		l.Error("RankInk error", err)
		return err
	}
	return nil
}

func RankHotTag(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: time.Second * 5,
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var activities *RankActivities
	l := workflow.GetLogger(ctx)
	err := workflow.ExecuteActivity(ctx, activities.RankTag, 100).Get(ctx, nil)
	if err != nil {
		l.Error("RankTag error", err)
		return err
	}
	return nil
}
