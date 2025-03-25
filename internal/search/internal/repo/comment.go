package repo

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/search/internal/domain"
	"github.com/KNICEX/InkFlow/internal/search/internal/repo/dao"
	"github.com/samber/lo"
)

type CommentRepo interface {
	Search(ctx context.Context, query string, offset, limit int) ([]domain.Comment, error)
	InputComment(ctx context.Context, comments []domain.Comment) error
	DeleteComment(ctx context.Context, commentId int64) error
	DeleteByBiz(ctx context.Context, biz string, bizId int64) error
}

type commentRepo struct {
	userParser
	dao     dao.CommentDAO
	userDAO dao.UserDAO
}

func NewCommentRepo(dao dao.CommentDAO, userDAO dao.UserDAO) CommentRepo {
	return &commentRepo{
		dao:     dao,
		userDAO: userDAO,
	}
}

func (repo *commentRepo) Search(ctx context.Context, query string, offset, limit int) ([]domain.Comment, error) {
	commentList, err := repo.dao.Search(ctx, query, offset, limit)
	if err != nil || len(commentList) == 0 {
		return nil, err
	}
	commentatorIds := lo.Map(commentList, func(item dao.Comment, index int) int64 {
		return item.CommentatorId
	})
	commentators, err := repo.userDAO.SearchByIds(ctx, commentatorIds)
	return lo.Map(commentList, func(item dao.Comment, index int) domain.Comment {
		comment := repo.entityToDomain(item)
		if commentator, ok := commentators[item.CommentatorId]; ok {
			comment.Commentator = repo.userParser.entityToDomain(commentator)
		}
		return comment
	}), nil
}

func (repo *commentRepo) InputComment(ctx context.Context, comments []domain.Comment) error {
	err := repo.dao.Input(ctx, lo.Map(comments, func(item domain.Comment, index int) dao.Comment {
		return repo.domainToEntity(item)
	}))
	if err != nil {
		return err
	}
	return nil
}

func (repo *commentRepo) DeleteComment(ctx context.Context, commentId int64) error {
	if err := repo.dao.DeleteByIds(ctx, []int64{commentId}); err != nil {
		return err
	}
	// 删除子评论
	return repo.dao.DeleteChildComments(ctx, commentId)
}

func (repo *commentRepo) DeleteByBiz(ctx context.Context, biz string, bizId int64) error {
	return repo.dao.DeleteByBiz(ctx, biz, bizId)
}

func (repo *commentRepo) domainToEntity(c domain.Comment) dao.Comment {
	return dao.Comment{
		Id:            c.Id,
		Biz:           c.Biz,
		BizId:         c.BizId,
		ParentId:      c.ParentId,
		RootId:        c.RootId,
		CommentatorId: c.Commentator.Id,
		Content:       c.Content,
		CreatedAt:     c.CreatedAt,
	}
}

func (repo *commentRepo) entityToDomain(c dao.Comment) domain.Comment {
	return domain.Comment{
		Id:       c.Id,
		Biz:      c.Biz,
		BizId:    c.BizId,
		ParentId: c.ParentId,
		RootId:   c.RootId,
		Commentator: domain.User{
			Id: c.CommentatorId,
		},
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
	}
}
