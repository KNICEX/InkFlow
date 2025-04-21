package schedule

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func ReviewRetryFail(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var activities *ReviewFailoverActivity
	l := workflow.GetLogger(ctx)
	err := workflow.ExecuteActivity(ctx, activities.RetryFail).Get(ctx, nil)
	if err != nil {
		l.Error("retry review fail error", err)
		return err
	}
	return nil
}
