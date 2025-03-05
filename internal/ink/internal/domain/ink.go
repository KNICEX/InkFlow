package domain

import "time"

type Ink struct {
	Id          int64
	Title       string
	Summary     string
	ContentHtml string
	// 可能引入块编辑器
	ContentMeta string
	Status      InkStatus
	Author      Author
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InkStatus int

type Author struct {
	Id      int64
	Name    string
	Account string
}
