package dao

import (
	"context"
	_ "embed"
	"github.com/elastic/go-elasticsearch/v8"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strings"
	"time"
)

var (
	//go:embed user_index.json
	userIndexMapping string
	//go:embed follow_index.json
	followIndexMapping string
	//go:embed ink_index.json
	inkIndexMapping string
)

const (
	userIndexName   = "user_index"
	followIndexName = "follow_index"
	inkIndexName    = "ink_index"
)

func InitEs(client *elasticsearch.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	eg := errgroup.Group{}
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, userIndexName, userIndexMapping)
	})
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, followIndexName, "")
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
