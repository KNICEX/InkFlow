package web

import "github.com/KNICEX/InkFlow/internal/search"

type SearchReq struct {
	Keyword string `json:"keyword" form:"keyword" binding:"required"`
	Type    string `json:"type" form:"type"`
	Offset  int    `json:"offset" form:"offset"`
	Limit   int    `json:"limit" form:"limit"`
}

func searchInkToInkVO(ink search.Ink) InkVO {
	return InkVO{
		Id:          ink.Id,
		Title:       ink.Title,
		Cover:       ink.Cover,
		Author:      searchUserToUserVO(ink.Author),
		ContentHtml: ink.Content,
		CreatedAt:   ink.CreatedAt,
		UpdatedAt:   ink.UpdatedAt,
		Tags:        ink.Tags,
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

func searchCommentToCommentVO(comment search.Comment) CommentVO {
	return CommentVO{
		Id:          comment.Id,
		Biz:         comment.Biz,
		BizId:       comment.BizId,
		Commentator: searchUserToUserVO(comment.Commentator),
		Payload: Payload{
			Content: comment.Content,
			Images:  comment.Images,
		},
		CreatedAt: comment.CreatedAt,
	}
}
