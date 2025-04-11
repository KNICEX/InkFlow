package ai

import (
	"github.com/KNICEX/InkFlow/internal/ai/internal/service"
	"github.com/KNICEX/InkFlow/internal/ai/internal/service/llm"
	"github.com/google/generative-ai-go/genai"
)

func InitLLMService(cli []*genai.Client) LLMService {
	svcs := make([]LLMService, 0, len(cli))
	for _, c := range cli {
		svcs = append(svcs, llm.NewGeminiService(c))
	}
	return service.NewFailoverService(svcs)
}
