package dao

import (
	"context"
	"github.com/KNICEX/InkFlow/pkg/mapstructurex"
	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
	"strconv"
	"time"
)

type Ink struct {
	Id        int64     `json:"id" mapstructure:"id"`
	Title     string    `json:"title" mapstructure:"title"`
	Cover     string    `json:"cover" mapstructure:"cover"`
	AuthorId  int64     `json:"author_id" mapstructure:"author_id"`
	Content   string    `json:"content" mapstructure:"content"`
	Tags      []string  `json:"tags" mapstructure:"tags"`
	AiTags    []string  `json:"ai_tags" mapstructure:"ai_tags"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at"`
}

type InkDAO interface {
	Search(ctx context.Context, key string, offset int, limit int) ([]Ink, error)
	InputInk(ctx context.Context, inks []Ink) error
	DeleteInk(ctx context.Context, ids []int64) error
}

type MeiliInkDAO struct {
	cli meilisearch.ServiceManager
}

func NewMeiliInkDAO(cli meilisearch.ServiceManager) InkDAO {
	return &MeiliInkDAO{cli: cli}
}

func (m *MeiliInkDAO) Search(ctx context.Context, key string, offset int, limit int) ([]Ink, error) {
	res, err := m.cli.Index(inkIndexName).SearchWithContext(ctx, key, &meilisearch.SearchRequest{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, err
	}
	if len(res.Hits) == 0 {
		return nil, nil
	}
	var inks []Ink
	err = mapstructurex.Decode(res.Hits, &inks)
	return inks, err
}

func (m *MeiliInkDAO) InputInk(ctx context.Context, inks []Ink) error {
	_, err := m.cli.Index(inkIndexName).AddDocumentsWithContext(ctx, inks)
	return err
}

func (m *MeiliInkDAO) DeleteInk(ctx context.Context, ids []int64) error {
	_, err := m.cli.Index(inkIndexName).DeleteDocumentsWithContext(ctx, lo.Map(ids, func(item int64, index int) string {
		return strconv.FormatInt(item, 10)
	}))
	return err
}
