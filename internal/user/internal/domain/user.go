package domain

import "time"

type User struct {
	Id int64
	// 账号名, 全局唯一， 以@开头展示
	Account string
	// 用户名，可重复
	Username string
	Email    string
	Phone    string
	Password string
	Link     []string
	AboutMe  string
	Level    int

	// oauth2
	GithubInfo GithubInfo

	Birthday  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
