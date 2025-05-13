package web

import (
	"time"

	"github.com/KNICEX/InkFlow/internal/comment"
)

type CommentVO struct {
	Id          int64        `json:"id,string"`
	Biz         string       `json:"biz"`
	BizId       int64        `json:"bizId,string"`
	Commentator UserVO       `json:"commentator"`
	IsAuthor    bool         `json:"isAuthor"`
	Payload     Payload      `json:"payload"`
	Parent      *CommentVO   `json:"parent"`
	Root        *CommentVO   `json:"root"`
	Stats       CommentStats `json:"stats"`
	Children    []CommentVO  `json:"children"`
	CreatedAt   time.Time    `json:"createdAt"`
}

type CommentStats struct {
	ReplyCnt int64 `json:"replyCnt"`
	LikeCnt  int64 `json:"likeCnt"`
	ViewCnt  int64 `json:"viewCnt"`
	Liked    bool  `json:"liked"`
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
		Stats:     commentStatsToVO(com.Stats),
		CreatedAt: com.CreatedAt,
	}
}

func commentStatsToVO(stats comment.Stats) CommentStats {
	return CommentStats{
		ReplyCnt: stats.ReplyCnt,
		LikeCnt:  stats.LikeCnt,
		Liked:    stats.Liked,
	}
}

type BizCommentReq struct {
	Biz   string `json:"biz" form:"biz"`
	BizId int64  `json:"bizId,string" form:"bizId"`
	MaxId int64  `json:"maxId,string" form:"maxId"`
	Limit int    `json:"limit" form:"limit"`
}

type ChildCommentReq struct {
	RootId int64 `json:"rootId,string" form:"rootId"`
	MaxId  int64 `json:"maxId,string" form:"maxId"`
	Limit  int   `json:"limit" form:"limit" binding:"required"`
}

type PostReplyReq struct {
	Biz      string  `json:"biz" binding:"required"`
	BizId    int64   `json:"bizId,string" binding:"required"`
	RootId   int64   `json:"rootId,string"`
	ParentId int64   `json:"parentId,string"`
	Payload  Payload `json:"payload" binding:"required"`
}
