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
	ContentMeta string   `json:"contentMeta"`
	Tags        []string `json:"tags"`
}

type InkDetailResp struct {
	InkBaseInfo
	Author      UserProfile   `json:"author"`
	Interactive InteractiveVO `json:"interactive"`
}

type InkCategory struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type InkBaseInfo struct {
	Id          int64       `json:"id"`
	Title       string      `json:"title"`
	Cover       string      `json:"cover"`
	Summary     string      `json:"summary"`
	Category    InkCategory `json:"category"`
	ContentType int         `json:"contentType"`
	Tags        []string    `json:"tags"`
	ContentHtml string      `json:"contentHtml"`
	ContentMeta string      `json:"contentMeta"`
	Status      int         `json:"status"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

func InkBaseInfoFromDomain(i ink.Ink) InkBaseInfo {
	return InkBaseInfo{
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
	AuthorId int64 `json:"authorId" binding:"required,max=100"`
	Category int64 `json:"category"`
	Offset   int   `json:"offset"`
	Limit    int   `json:"limit"`
}

type ListSelfReq struct {
	Category int64 `json:"category"`
	Offset   int   `json:"offset"`
	Limit    int   `json:"limit"`
}

type ListDraftReq struct {
	Category int64 `json:"category"`
	Offset   int   `json:"offset"`
	Limit    int   `json:"limit"`
}

type ListMaxIdReq struct {
	MaxId int64 `json:"maxId"`
	Limit int   `json:"limit"`
}
