package gemini

import (
	"context"
	"errors"
	"github.com/KNICEX/InkFlow/internal/ai/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ai/internal/service/llm"
	"github.com/google/generative-ai-go/genai"
	"strings"
)

type Service struct {
	key    string
	client *genai.Client
	model  *genai.GenerativeModel
}

type Session struct {
	session *genai.ChatSession
}

func (c *Session) Ask(ctx context.Context, question string) (domain.Resp, error) {
	resp, err := c.session.SendMessage(ctx, genai.Text(question))
	if err != nil {
		return domain.Resp{}, err
	}
	res, token := parseResponse(resp)
	return domain.Resp{
		Content: res,
		Token:   token,
	}, nil
}

func (c *Session) Close() error {
	return nil
}

type Option func(*Service)

func NewGeminiService(client *genai.Client, opts ...Option) llm.Service {
	svc := &Service{
		client: client,
		model:  client.GenerativeModel("gemini-2.0-flash"),
	}

	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func WithPreset(prompt string) Option {
	return func(s *Service) {
		s.model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{
				genai.Text(prompt),
			},
		}
	}
}

func WithTemperature(temp float32) Option {
	return func(s *Service) {
		s.model.SetTemperature(temp)
	}
}

func WithFlash2() Option {
	return func(s *Service) {
		s.model = s.client.GenerativeModel("gemini-2.0-flash")
	}
}

func (svc *Service) parseResp(resp *genai.GenerateContentResponse) (string, error) {
	if len(resp.Candidates) == 0 {
		return "", errors.New("no response")
	}
	text := resp.Candidates[0].Content.Parts[0].(genai.Text)
	return string(text), nil
}

func (svc *Service) AskOnce(ctx context.Context, msg string) (domain.Resp, error) {
	resp, err := svc.model.GenerateContent(ctx, genai.Text(msg))
	if err != nil {
		return domain.Resp{}, err
	}
	res, token := parseResponse(resp)
	return domain.Resp{
		Content: res,
		Token:   token,
	}, nil
}

func (svc *Service) BeginChat(ctx context.Context) (llm.Session, error) {
	session := svc.model.StartChat()
	return &Session{
		session: session,
	}, nil

}

func parseResponse(resp *genai.GenerateContentResponse) (string, int64) {
	var resStr strings.Builder
	if resp.Candidates != nil && len(resp.Candidates) > 0 {
		for i, part := range resp.Candidates[0].Content.Parts {
			if part == nil {
				continue
			}
			if text, ok := part.(genai.Text); ok {
				if i > 0 {
					resStr.WriteString("\n")
				}
				resStr.WriteString(string(text))
			} else {
				return "", 0
			}
		}
	}
	return resStr.String(), int64(resp.UsageMetadata.TotalTokenCount)
}
