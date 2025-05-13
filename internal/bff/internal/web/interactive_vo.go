package web

import (
	"time"

	"github.com/KNICEX/InkFlow/internal/interactive"
)

type InteractiveVO struct {
	Biz         string `json:"biz"`
	BizId       int64  `json:"bizId,string"`
	ViewCnt     int64  `json:"viewCnt"`
	LikeCnt     int64  `json:"likeCnt"`
	CommentCnt  int64  `json:"commentCnt"`
	ShareCnt    int64  `json:"shareCnt"`
	FavoriteCnt int64  `json:"favoriteCnt"`
	Liked       bool   `json:"liked"`
	Favorited   bool   `json:"favorited"`
	Shared      bool   `json:"shared"`
}

func intrToVo(i interactive.Interactive) InteractiveVO {
	return InteractiveVO{
		Biz:         i.Biz,
		BizId:       i.BizId,
		ViewCnt:     i.ViewCnt,
		LikeCnt:     i.LikeCnt,
		FavoriteCnt: i.FavoriteCnt,
		Liked:       i.Liked,
		Favorited:   i.Favorited,
	}
}

type FavoriteVO struct {
	Fid       int64     `json:"fid,string"`
	Name      string    `json:"name"`
	Biz       string    `json:"biz"`
	Private   bool      `json:"private"`
	CreatedAt time.Time `json:"createdAt"`
}

func favoriteToVO(f interactive.Favorite) FavoriteVO {
	return FavoriteVO{
		Fid:       f.Id,
		Name:      f.Name,
		Biz:       f.Biz,
		Private:   f.Private,
		CreatedAt: f.CreatedAt,
	}
}

type CreateFavReq struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
	Biz     string `json:"biz"`
}
