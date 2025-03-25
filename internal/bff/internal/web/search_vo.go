package web

import "github.com/KNICEX/InkFlow/internal/search"

type SearchReq struct {
	Keyword string `json:"keyword"`
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
}

func searchInkToInkVO(ink search.Ink) InkVO {
	return InkVO{
		InkBaseVO: InkBaseVO{
			Id:          ink.Id,
			Title:       ink.Title,
			Cover:       ink.Cover,
			ContentHtml: ink.Content,
			CreatedAt:   ink.CreatedAt,
			UpdatedAt:   ink.UpdatedAt,
			Tags:        ink.Tags,
		},
		Author: searchUserToUserVO(ink.Author),
	}
}

func searchUserToUserVO(user search.User) UserVO {
	return UserVO{
		Id:        user.Id,
		Username:  user.Username,
		Account:   user.Account,
		AboutMe:   user.AboutMe,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
	}
}
