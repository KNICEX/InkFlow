package web

import "time"

type CommentVO struct {
	Id          int64       `json:"id"`
	Biz         string      `json:"biz"`
	BizId       int64       `json:"bizId"`
	Commentator UserVO      `json:"commentator"`
	Content     string      `json:"content"`
	Images      []string    `json:"images"`
	Parent      *CommentVO  `json:"parent"`
	Root        *CommentVO  `json:"root"`
	Children    []CommentVO `json:"children"`
	CreatedAt   time.Time   `json:"createdAt"`
}
