package domain

import "time"

const (
	BizInk     = "ink"
	BizComment = "comment"
)

type Interactive struct {
	Biz         string
	BizId       int64
	ViewCnt     int64
	LikeCnt     int64
	ReplyCnt    int64
	FavoriteCnt int64

	Liked     bool
	Favorited bool
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

type FavoriteRecord struct {
	Biz       string
	BizId     int64
	Fid       int64
	UserId    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
