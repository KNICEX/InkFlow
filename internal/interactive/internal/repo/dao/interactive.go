package dao

import "time"

type UserRead struct {
	Id        int64
	UserId    int64  `gorm:"uniqueIndex:userId_biz_id_idx"`
	BizId     int64  `gorm:"uniqueIndex:userId_biz_id_idx"`
	Biz       string `gorm:"type:varchar(64);uniqueIndex:userId_biz_id_idx"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"index"`
}

type UserLike struct {
	Id        int64
	UserId    int64  `gorm:"uniqueIndex:userId_biz_id_idx"`
	BizId     int64  `gorm:"uniqueIndex:userId_biz_id_idx"`
	Biz       string `gorm:"type:varchar(64);uniqueIndex:userId_biz_id_idx"`
	Status    int
	UpdatedAt time.Time `gorm:"index"`
	CreatedAt time.Time
}

type UserUnlike struct {
	Id        int64
	UserId    int64  `gorm:"uniqueIndex:userId_biz_id_idx"`
	BizId     int64  `gorm:"uniqueIndex:userId_biz_id_idx"`
	Biz       string `gorm:"type:varchar(64);uniqueIndex:userId_biz_id_idx"`
	Status    int
	UpdatedAt time.Time `gorm:"index"`
	CreatedAt time.Time
}

// UserCollection TODO 考虑支持多个收藏夹
type UserCollection struct {
	Id           int64
	UserId       int64  `gorm:"uniqueIndex:userId_biz_id_idx"`
	BizId        int64  `gorm:"uniqueIndex:userId_biz_id_idx"`
	Biz          string `gorm:"type:varchar(64);uniqueIndex:userId_biz_id_idx"`
	CollectionId int64  `gorm:"index;uniqueIndex:userId_biz_id_idx"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Interactive struct {
	Id    int64
	BizId int64  `gorm:"uniqueIndex:biz_type_idx"`
	Biz   string `gorm:"type:varchar(64);uniqueIndex:biz_type_idx"`

	ReadCnt    int64
	LikeCnt    int64
	UnlikeCnt  int64
	CollectCnt int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
