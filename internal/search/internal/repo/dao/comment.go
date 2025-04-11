package dao

import (
	"context"
	"fmt"
	"github.com/KNICEX/InkFlow/pkg/mapstructurex"
	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
	"strconv"
	"time"
)

type Comment struct {
	Id            int64     `json:"id" mapstructure:"id"`
	Biz           string    `json:"biz" mapstructure:"biz"`
	BizId         int64     `json:"biz_id" mapstructure:"biz_id"`
	RootId        int64     `json:"root_id" mapstructure:"root_id"`
	ParentId      int64     `json:"parent_id" mapstructure:"parent_id"`
	CommentatorId int64     `json:"commentator_id" mapstructure:"commentator_id"`
	Content       string    `json:"content" mapstructure:"content"`
	Images        string    `json:"images" mapstructure:"images"`
	CreatedAt     time.Time `json:"created_at" mapstructure:"created_at"`
}

type CommentDAO interface {
	Search(ctx context.Context, keyword string, offset, limit int) ([]Comment, error)
	Input(ctx context.Context, comments []Comment) error
	DeleteByIds(ctx context.Context, ids []int64) error
	DeleteChildComments(ctx context.Context, id int64) error
	DeleteByBiz(ctx context.Context, biz string, bizId int64) error
}

type MeiliCommentDAO struct {
	cli meilisearch.ServiceManager
}

func NewMeiliCommentDAO(cli meilisearch.ServiceManager) CommentDAO {
	return &MeiliCommentDAO{cli: cli}
}

func (m *MeiliCommentDAO) Search(ctx context.Context, keyword string, offset, limit int) ([]Comment, error) {
	res, err := m.cli.Index(commentIndexName).SearchWithContext(ctx, keyword, &meilisearch.SearchRequest{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, err
	}
	if len(res.Hits) == 0 {
		return nil, nil
	}
	var comments []Comment
	err = mapstructurex.Decode(res.Hits, &comments)
	return comments, err
}

func (m *MeiliCommentDAO) Input(ctx context.Context, comments []Comment) error {
	_, err := m.cli.Index(commentIndexName).AddDocumentsWithContext(ctx, comments)
	return err
}

func (m *MeiliCommentDAO) DeleteByIds(ctx context.Context, ids []int64) error {
	_, err := m.cli.Index(commentIndexName).DeleteDocumentsWithContext(ctx, lo.Map(ids, func(item int64, index int) string {
		return strconv.FormatInt(item, 10)
	}))
	return err
}

func (m *MeiliCommentDAO) DeleteChildComments(ctx context.Context, id int64) error {
	_, err := m.cli.Index(commentIndexName).DeleteDocumentsByFilterWithContext(ctx,
		fmt.Sprintf("root_id=%d OR parent_id=%d", id, id))
	return err
}

func (m *MeiliCommentDAO) DeleteByBiz(ctx context.Context, biz string, bizId int64) error {
	_, err := m.cli.Index(commentIndexName).DeleteDocumentsByFilterWithContext(ctx,
		fmt.Sprintf("biz=%s AND biz_id=%d", biz, bizId))
	return err
}
