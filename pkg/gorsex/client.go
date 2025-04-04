package gorsex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	client "github.com/gorse-io/gorse-go"
	"io"
	"log"
	"net/http"
	"strings"
)

type Client struct {
	*client.GorseClient
	entrypoint string
	apiKey     string
	httpClient http.Client
}

func NewClient(entrypoint, apiKey string) *Client {
	return &Client{
		GorseClient: client.NewGorseClient(entrypoint, apiKey),
		entrypoint:  entrypoint,
		apiKey:      apiKey,
		httpClient:  http.Client{},
	}
}

func (c *Client) DeleteFeedback(ctx context.Context, feedback client.Feedback) (client.RowAffected, error) {
	return request[client.RowAffected, any](ctx, c, http.MethodDelete, fmt.Sprintf("%s/api/feedback/%s/%s/%s", c.entrypoint, feedback.FeedbackType, feedback.UserId, feedback.ItemId), nil)
}

func request[Response any, Body any](ctx context.Context, c *Client, method, url string, body Body) (result Response, err error) {
	bodyByte, marshalErr := json.Marshal(body)
	if marshalErr != nil {
		return result, marshalErr
	}
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, method, url, strings.NewReader(string(bodyByte)))
	if err != nil {
		return result, err
	}
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return result, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		return result, errors.New(buf.String())
	}
	err = json.Unmarshal([]byte(buf.String()), &result)
	if err != nil {
		return result, err
	}
	return result, err
}
