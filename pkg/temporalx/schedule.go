package temporalx

import (
	"context"
	"errors"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

func UpsertSchedule(ctx context.Context, c client.Client, option client.ScheduleOptions) error {
	handler := c.ScheduleClient().GetHandle(ctx, option.ID)
	desc, err := handler.Describe(ctx)
	if err != nil {
		var notFound *serviceerror.NotFound
		if errors.As(err, &notFound) {
			// 不存在，创建
			handler, err = c.ScheduleClient().Create(ctx, option)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	// 存在，更新
	err = handler.Update(ctx, client.ScheduleUpdateOptions{
		DoUpdate: func(input client.ScheduleUpdateInput) (*client.ScheduleUpdate, error) {
			return &client.ScheduleUpdate{
				Schedule: &client.Schedule{
					Spec:   &option.Spec,
					Action: option.Action,
					Policy: &client.SchedulePolicies{
						Overlap:        option.Overlap,
						CatchupWindow:  option.CatchupWindow,
						PauseOnFailure: option.PauseOnFailure,
					},
					State: input.Description.Schedule.State,
				},
			}, nil
		},
	})
	if err != nil {
		return err
	}

	if desc.Schedule.State.Paused {
		// 如果是暂停状态，恢复
		return handler.Unpause(ctx, client.ScheduleUnpauseOptions{})
	}
	return nil
}
