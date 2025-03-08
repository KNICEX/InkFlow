package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiService struct {
	key    string
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiService(key string) (*GeminiService, error) {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	return &GeminiService{
		key:    key,
		client: client,
		model:  client.GenerativeModel("gemini-1.5-flash"),
	}, nil
}

func (g *GeminiService) parseResp(resp *genai.GenerateContentResponse) (string, error) {
	if len(resp.Candidates) == 0 {
		return "", errors.New("no response")
	}
	text := resp.Candidates[0].Content.Parts[0].(genai.Text)
	return string(text), nil
}

func (g *GeminiService) AskOnce(ctx context.Context, msg string) (Resp, error) {
	resp, err := g.model.GenerateContent(ctx, genai.Text(msg))
	if err != nil {
		return Resp{}, err
	}
	if len(resp.Candidates) == 0 {
		return Resp{}, errors.New("no response")
	}
	content, err := g.parseResp(resp)
	if err != nil {
		return Resp{}, err
	}
	return Resp{Content: content}, nil
}

func (g *GeminiService) BeginChat(ctx context.Context, msg string) <-chan Context {
	ch := make(chan Context, 1)
	chat := g.model.StartChat()
	var sendFunc func(ctx context.Context, msg string) error
	sendFunc = func(ctx context.Context, msg string) error {
		resp, err := chat.SendMessage(ctx, genai.Text(msg))
		if err != nil {
			return err
		}
		content, err := g.parseResp(resp)
		if err != nil {
			return err
		}
		ch <- Context{
			Resp: Resp{Content: content},
			Ask:  sendFunc,
		}
		return nil
	}

	go func() {
		err := sendFunc(ctx, msg)
		if err != nil {
			ch <- Context{
				Resp: Resp{Content: err.Error()},
				Err:  err,
				Ask:  sendFunc,
			}
		}
	}()

	return ch
}

func Example() error {
	svc, err := NewGeminiService("your-api-key")
	if err != nil {
		return err
	}

	ch := svc.BeginChat(context.Background(), "hello")
	for ctx := range ch {
		if ctx.Err != nil {
			return ctx.Err
		}
		fmt.Println(ctx.Resp.Content)
		if err := ctx.Ask(context.Background(), "world"); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
