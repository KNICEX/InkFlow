package dao

import (
	"context"
	_ "embed"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/meilisearch/meilisearch-go"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strings"
	"time"
)

var (
	//go:embed user_index.json
	userIndexMapping string
	//go:embed ink_index.json
	inkIndexMapping string
)

const (
	userIndexName    = "user_index"
	commentIndexName = "comment_index"
	inkIndexName     = "ink_index"
)

func InitEs(client *elasticsearch.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	eg := errgroup.Group{}
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, userIndexName, userIndexMapping)
	})
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, inkIndexName, inkIndexMapping)
	})
	return eg.Wait()
}

func tryCreateIndex(ctx context.Context,
	client *elasticsearch.Client,
	indexName string,
	indexMap string) error {
	resp, err := client.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	_, err = client.Indices.Create(indexName,
		client.Indices.Create.WithContext(ctx),
		client.Indices.Create.WithBody(strings.NewReader(indexMap)),
	)
	if err != nil {
		return err
	}
	return nil
}

func InitMeili(cli meilisearch.ServiceManager) error {
	_, err := cli.CreateIndex(&meilisearch.IndexConfig{
		Uid:        userIndexName,
		PrimaryKey: "id",
	})
	if err != nil {
		return err
	}
	_, err = cli.CreateIndex(&meilisearch.IndexConfig{
		Uid:        inkIndexName,
		PrimaryKey: "id",
	})
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 3)
	return initIndexSetting(cli)
}

func initIndexSetting(cli meilisearch.ServiceManager) error {
	// 设置可搜索属性
	_, err := cli.Index(userIndexName).UpdateSearchableAttributes(&[]string{
		"account",
		"username",
		"about_me",
	})
	if err != nil {
		return err
	}
	// 设置可过滤属性, 通过类sql语法进行过滤
	_, err = cli.Index(userIndexName).UpdateFilterableAttributes(&[]string{
		"id",
		"account",
	})
	if err != nil {
		return err
	}
	_, err = cli.Index(inkIndexName).UpdateSearchableAttributes(&[]string{
		"title",
		"content",
		"tags",
		"ai_tags",
	})
	if err != nil {
		return err
	}
	_, err = cli.Index(inkIndexName).UpdateFilterableAttributes(&[]string{
		"id",
		"tags",
	})
	_, err = cli.Index(commentIndexName).UpdateFilterableAttributes(&[]string{
		"content",
	})
	return err
}
