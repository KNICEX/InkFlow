package domain

import (
	"strings"
	"time"
)

type User struct {
	Id int64
	// 账号名, 全局唯一， 以@开头展示
	Account string
	// 用户名，可重复
	Username string
	Email    string
	Phone    string
	Password string
	Links    Links
	AboutMe  string
	Level    int

	// oauth2
	GithubInfo GithubInfo

	Birthday  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Links []string

func (links Links) ToString() string {
	sb := strings.Builder{}
	for _, link := range links {
		sb.WriteString(link)
		sb.WriteString(",")
	}
	return sb.String()
}

func LinksFromString(s string) Links {
	return strings.Split(s, ",")
}
