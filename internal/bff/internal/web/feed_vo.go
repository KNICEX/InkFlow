package web

import "github.com/KNICEX/InkFlow/internal/feed"

type FeedFollowReq struct {
	MaxId     int64 `json:"maxId,string" form:"maxId"`
	Timestamp int64 `json:"timestamp,string" form:"timestamp"`
	Limit     int   `json:"limit" form:"limit"`
}

type FeedRecommendReq struct {
	Offset int `json:"offset" form:"offset"`
	Limit  int `json:"limit" form:"limit"`
}

func feedInkToVO(feedInk feed.Ink) InkVO {
	return InkVO{
		Id: feedInk.InkId,
		Author: UserVO{
			Id: feedInk.AuthorId,
		},
		ContentHtml: feedInk.Abstract,
		Cover:       feedInk.Cover,
		Title:       feedInk.Title,
		CreatedAt:   feedInk.CreatedAt,
		UpdatedAt:   feedInk.CreatedAt,
	}
}
