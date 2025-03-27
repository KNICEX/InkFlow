package web

import (
	"github.com/KNICEX/InkFlow/internal/comment"
	"time"
)

type CommentVO struct {
	Id          int64       `json:"id"`
	Biz         string      `json:"biz"`
	BizId       int64       `json:"bizId"`
	Commentator UserVO      `json:"commentator"`
	IsAuthor    bool        `json:"isAuthor"`
	Payload     Payload     `json:"payload"`
	Parent      *CommentVO  `json:"parent"`
	Root        *CommentVO  `json:"root"`
	Children    []CommentVO `json:"children"`
	CreatedAt   time.Time   `json:"createdAt"`
}

type Payload struct {
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

func commentToVO(com comment.Comment) CommentVO {
	var parent, root *CommentVO
	if com.Parent != nil {
		vo := commentToVO(*com.Parent)
		parent = &vo
	}
	if com.Root != nil {
		vo := commentToVO(*com.Root)
		root = &vo
	}
	var children []CommentVO
	if len(com.Children) > 0 {
		children = make([]CommentVO, len(com.Children))
		for i, child := range com.Children {
			vo := commentToVO(child)
			children[i] = vo
		}
	}
	return CommentVO{
		Id:    com.Id,
		Biz:   com.Biz,
		BizId: com.BizId,
		Commentator: UserVO{
			Id: com.Commentator.Id,
		},
		IsAuthor: com.Commentator.IsAuthor,
		Payload: Payload{
			Content: com.Payload.Content,
			Images:  com.Payload.Images,
		},
		Parent:    parent,
		Root:      root,
		Children:  children,
		CreatedAt: com.CreatedAt,
	}
}
