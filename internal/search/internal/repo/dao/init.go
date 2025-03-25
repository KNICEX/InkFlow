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
	return nil
}
