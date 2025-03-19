package dao

import (
	"context"
	"github.com/KNICEX/InkFlow/pkg/elasticx"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"strconv"
	"time"
)

type User struct {
	Id          int64     `json:"id"`
	Email       string    `json:"email"`
	Account     string    `json:"account"`
	Username    string    `json:"username"`
	FollowerCnt int64     `json:"follower_cnt"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserDAO interface {
	SearchByUsername(ctx context.Context, keywords []string) ([]User, error)
	SearchByAccount(ctx context.Context, keywords []string) ([]User, error)
	InputUser(ctx context.Context, user User) error
	BatchInputUser(ctx context.Context, users []User) error
}

type ElasticUserDAO struct {
	client      *elasticsearch.Client
	typedClient *elasticsearch.TypedClient
}

func NewElasticUserDAO(client *elasticsearch.Client) UserDAO {
	typedApi := typedapi.New(client)

	return ElasticUserDAO{
		client: client,
		typedClient: &elasticsearch.TypedClient{
			BaseClient: client.BaseClient,
			API:        typedApi,
		},
	}
}

func (dao ElasticUserDAO) SearchByUsername(ctx context.Context, keywords []string) ([]User, error) {
	must := make([]types.Query, 0, len(keywords))
	for _, keyword := range keywords {
		must = append(must, types.Query{
			Match: map[string]types.MatchQuery{
				"username": {
					Query: keyword,
				},
			},
		})
	}
	resp, err := dao.typedClient.Search().
		Index(userIndexName).
		Request(&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Must: must,
				},
			},
		}).Do(ctx)
	if err != nil {
		return nil, err
	}
	return elasticx.MarshalResp[User](resp)
}

func (dao ElasticUserDAO) SearchByAccount(ctx context.Context, keywords []string) ([]User, error) {
	//TODO implement me
	panic("implement me")
}

func (dao ElasticUserDAO) InputUser(ctx context.Context, user User) error {
	_, err := dao.typedClient.Index(userIndexName).
		Id(strconv.FormatInt(user.Id, 10)).
		Request(user).Do(ctx)

	return err
}

func (dao ElasticUserDAO) BatchInputUser(ctx context.Context, users []User) error {
	panic("")
}
