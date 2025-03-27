package domain

import "time"

type Comment struct {
	Id    int64
	Biz   string
	BizId int64

	Root        *Comment
	Parent      *Comment
	Commentator Commentator
	Children    []Comment
	Payload     Payload
	Stats       CommentStats
	CreatedAt   time.Time
}

type Commentator struct {
	Id       int64
	IsAuthor bool
}

type CommentStats struct {
	LikeCnt  int64
	ReplyCnt int64
	Liked    bool
}

type Payload struct {
	Content string
	Images  []string
}
