package elasticx

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi"
)

func ClientToTypedClient(c *elasticsearch.Client) *elasticsearch.TypedClient {
	typedApi := typedapi.New(c)
	return &elasticsearch.TypedClient{
		BaseClient: c.BaseClient,
		API:        typedApi,
	}
}
