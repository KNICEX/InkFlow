package domain

import "time"

type User struct {
	Id        int64
	Account   string
	CreatedAt time.Time
}
