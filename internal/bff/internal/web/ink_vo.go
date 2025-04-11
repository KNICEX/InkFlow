package web

import (
	"github.com/KNICEX/InkFlow/internal/ink"
	"time"
)

type SaveInkReq struct {
	Id          int64    `json:"id"`
	Title       string   `json:"title" binding:"required,max=100"`
	Cover       string   `json:"cover"`
	Summary     string   `json:"summary"`
	ContentHtml string   `json:"contentHtml"`
	ContentMeta string   `json:"contentMeta" binding:"required"`
	Tags        []string `json:"tags"`
}

type InkVO struct {
	Id          int64         `json:"id"`
	Title       string        `json:"title"`
	Author      UserVO        `json:"author"`
	Cover       string        `json:"cover"`
	Summary     string        `json:"summary"`
	Category    InkCategory   `json:"category"`
	ContentType int           `json:"contentType"`
	Tags        []string      `json:"tags"`
	ContentHtml string        `json:"contentHtml"`
	ContentMeta string        `json:"contentMeta"`
	Status      int           `json:"status"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	Interactive InteractiveVO `json:"interactive"`
}

type InkCategory struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func inkToVO(i ink.Ink) InkVO {
	return InkVO{
		Id:          i.Id,
		Title:       i.Title,
		Summary:     i.Summary,
		Cover:       i.Cover,
		Tags:        i.Tags,
		Category:    InkCategory{Id: i.Category.Id},
		ContentType: i.ContentType.ToInt(),
		ContentHtml: i.ContentHtml,
		ContentMeta: i.ContentMeta,
		Status:      i.Status.ToInt(),
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

type LikeReq struct {
}

type ListReq struct {
	AuthorId int64 `json:"authorId" form:"authorId" binding:"required"`
	Category int64 `json:"category" form:"category"`
	Offset   int   `json:"offset" form:"offset"`
	Limit    int   `json:"limit" form:"limit" binding:"required,max=500"`
}

type ListSelfReq struct {
	Category int64 `json:"category" form:"category"`
	Offset   int   `json:"offset" form:"offset"`
	Limit    int   `json:"limit" form:"limit"`
}

type ListDraftReq struct {
	Category int64 `json:"category" form:"category"`
	Offset   int   `json:"offset" form:"offset"`
	Limit    int   `json:"limit" form:"limit"`
}

type ListMaxIdReq struct {
	MaxId int64 `json:"maxId" form:"maxId"`
	Limit int   `json:"limit" form:"limit"`
}

type ListFavoriteReq struct {
	Fid   int64 `json:"fid" form:"fid"`
	MaxId int64 `json:"maxId" form:"maxId"`
	Limit int   `json:"limit" form:"limit"`
}

type FavoriteReq struct {
	FavoriteId int64 `json:"favoriteId" from:"favoriteId"`
}
