package web

import "github.com/KNICEX/InkFlow/internal/interactive"

type InteractiveVO struct {
	Biz        string `json:"biz"`
	BizId      int64  `json:"bizId"`
	ViewCnt    int64  `json:"viewCnt"`
	LikeCnt    int64  `json:"likeCnt"`
	CommentCnt int64  `json:"commentCnt"`
	ShareCnt   int64  `json:"shareCnt"`
	CollectCnt int64  `json:"collectCnt"`
	Liked      bool   `json:"liked"`
	Favorited  bool   `json:"favorited"`
	Shared     bool   `json:"shared"`
}

func InteractiveVOFromDomain(i interactive.Interactive) InteractiveVO {
	return InteractiveVO{
		Biz:        i.Biz,
		BizId:      i.BizId,
		ViewCnt:    i.ViewCnt,
		LikeCnt:    i.LikeCnt,
		CollectCnt: i.CollectCnt,
		Liked:      i.Liked,
		Favorited:  i.Favorited,
	}
}
