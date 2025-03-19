package service

import (
	"context"
)

type SubscribeService interface {
	// Subscribe subscribes to a user who has been followed.
	Subscribe(ctx context.Context, uid, subUid int64) error
	Unsubscribe(ctx context.Context, uid, subUid int64) error
	SubscribeList(ctx context.Context, uid int64, offset, limit int) ([]int64, error)
}
