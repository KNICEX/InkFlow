package domain

import "time"

type Interactive struct {
	Biz        string
	BizId      int64
	ViewCnt    int64
	LikeCnt    int64
	CollectCnt int64

	Liked     bool
	Collected bool
}

type ViewRecord struct {
	Biz       string
	BizId     int64
	UserId    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LikeRecord struct {
	Biz       string
	BizId     int64
	UserId    int64
	CreatedAt time.Time
}
