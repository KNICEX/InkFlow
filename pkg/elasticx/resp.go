package elasticx

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
)

func MarshalResp[T any](resp *search.Response) ([]T, error) {
	res := make([]T, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var t T
		if err := json.Unmarshal(hit.Source_, &t); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}
