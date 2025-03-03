package dao

import (
	"context"
	"time"
)

type UserFollow struct {
	Id       int64
	Follower int64 `gorm:"uniqueIndex:follower_followee"`
	Followee int64 `gorm:"uniqueIndex:follower_followee"`

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}

// Block TODO 后续考虑支持
type Block struct {
	Id        int64
	Blocker   int64
	Blocked   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FollowStatistic struct {
	Id         int64
	UserId     int64 `gorm:"unique"`
	Followers  int64
	Followings int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRelationDAO interface {
	FollowList(ctx context.Context, uid int64, offset, limit int) ([]UserFollow, error)
	FollowDetail(ctx context.Context, follower, followee int64) (UserFollow, error)
	CreateFollowRelation(ctx context.Context, c UserFollow) error
	CntFollower(ctx context.Context, uid int64) (int64, error)
	CntFollowee(ctx context.Context, uid int64) (int64, error)
	FollowStatistic(ctx context.Context, uid int64) (FollowStatistic, error)
}
