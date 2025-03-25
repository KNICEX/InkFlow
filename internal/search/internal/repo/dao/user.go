package dao

import (
	"context"
	"github.com/meilisearch/meilisearch-go"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/lo"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Id        int64     `json:"id"`
	Avatar    string    `json:"avatar"`
	Account   string    `json:"account"`
	Username  string    `json:"username"`
	AboutMe   string    `json:"about_me"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserDAO interface {
	Search(ctx context.Context, query string, offset, limit int) ([]User, error)
	SearchByIds(ctx context.Context, ids []int64) (map[int64]User, error)
	InputUser(ctx context.Context, users []User) error
	DeleteUser(ctx context.Context, userIds []int64) error
}

type MeiliUserDAO struct {
	cli meilisearch.ServiceManager
}

func NewMeiliUserDAO(cli meilisearch.ServiceManager) UserDAO {
	return &MeiliUserDAO{cli: cli}
}

func (dao *MeiliUserDAO) Search(ctx context.Context, query string, offset, limit int) ([]User, error) {
	res, err := dao.cli.Index(userIndexName).SearchWithContext(ctx, query, &meilisearch.SearchRequest{
		Offset: int64(offset),
		Limit:  int64(limit),
	})
	if err != nil {
		return nil, err
	}
	if len(res.Hits) == 0 {
		return nil, nil
	}
	var users []User
	err = mapstructure.Decode(res.Hits, &users)
	return users, err
}

func (dao *MeiliUserDAO) SearchByIds(ctx context.Context, ids []int64) (map[int64]User, error) {
	idsStr := strings.Builder{}
	for i, id := range ids {
		if i > 0 {
			idsStr.WriteString(",")
		}
		idsStr.WriteString(strconv.FormatInt(id, 10))
	}
	res, err := dao.cli.Index(userIndexName).SearchWithContext(ctx, "", &meilisearch.SearchRequest{
		Filter: []string{"id IN [" + idsStr.String() + "]"},
	})
	if err != nil {
		return nil, err
	}
	if len(res.Hits) == 0 {
		return nil, nil
	}
	var users []User
	err = mapstructure.Decode(res.Hits, &users)
	if err != nil {
		return nil, err
	}
	userMap := make(map[int64]User)
	for _, user := range users {
		userMap[user.Id] = user
	}
	return userMap, nil
}

func (dao *MeiliUserDAO) InputUser(ctx context.Context, users []User) error {
	_, err := dao.cli.Index(userIndexName).AddDocumentsWithContext(ctx, users)
	return err
}

func (dao *MeiliUserDAO) DeleteUser(ctx context.Context, userIds []int64) error {
	_, err := dao.cli.Index(userIndexName).DeleteDocumentsWithContext(ctx, lo.Map(userIds, func(item int64, index int) string {
		return strconv.FormatInt(item, 10)
	}))
	return err
}
