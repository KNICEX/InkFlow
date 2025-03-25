package domain

import "time"

type User struct {
	Id        int64
	Avatar    string
	Account   string
	Username  string
	AboutMe   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
