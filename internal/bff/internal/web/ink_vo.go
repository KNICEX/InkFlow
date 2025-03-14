package web

import "github.com/KNICEX/InkFlow/internal/ink"

type SaveInkReq struct {
	Id          int64    `json:"id"`
	Title       string   `json:"title"`
	Cover       string   `json:"cover"`
	Summary     string   `json:"summary"`
	ContentHtml string   `json:"content_html"`
	ContentMeta string   `json:"content_meta"`
	Tags        []string `json:"tags"`
}

type PublishInkReq struct {
	Id int64 `json:"id"`
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
	ContentType int         `json:"content_type"`
	Tags        []string    `json:"tags"`
	ContentHtml string      `json:"content_html"`
	ContentMeta string      `json:"content_meta"`
	Status      int         `json:"status"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
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
		CreatedAt:   i.CreatedAt.String(),
		UpdatedAt:   i.UpdatedAt.String(),
	}
}

type LikeReq struct {
}

type ListReq struct {
	AuthorId int64 `json:"author_id" bind:"required,max=100"`
	Category int64 `json:"category"`
	Offset   int   `json:"offset"`
	Limit    int   `json:"limit"`
}

type ListDraftReq struct {
	Category int64 `json:"category"`
	Offset   int   `json:"offset"`
	Limit    int   `json:"limit"`
}
