package web

import "github.com/KNICEX/InkFlow/internal/interactive"

type InteractiveVO struct {
	Biz        string `json:"biz"`
	BizId      int64  `json:"biz_id"`
	ViewCnt    int64  `json:"view_cnt"`
	LikeCnt    int64  `json:"like_cnt"`
	CommentCnt int64  `json:"comment_cnt"`
	ShareCnt   int64  `json:"share_cnt"`
	CollectCnt int64  `json:"collect_cnt"`
	Liked      bool   `json:"liked"`
	Collected  bool   `json:"collected"`
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
		Collected:  i.Collected,
	}
}
