package domain

import "time"

type Favorite struct {
	Id        int64
	UserId    int64
	Name      string
	Biz       string
	Private   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
