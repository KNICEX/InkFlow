package schedule

import "go.temporal.io/sdk/workflow"

func RankHotInk(ctx workflow.Context, n int) error {
	var activities *RankActivities
	l := workflow.GetLogger(ctx)
	err := workflow.ExecuteActivity(ctx, activities.RankInk, n).Get(ctx, nil)
	if err != nil {
		l.Error("RankInk error", err)
		return err
	}
	return nil
}

func RankHotTag(ctx workflow.Context, n int) error {
	var activities *RankActivities
	l := workflow.GetLogger(ctx)
	err := workflow.ExecuteActivity(ctx, activities.RankTag, n).Get(ctx, nil)
	if err != nil {
		l.Error("RankTag error", err)
		return err
	}
	return nil
}
