package ai

import (
	"github.com/KNICEX/InkFlow/internal/ai/internal/service/llm"
	"github.com/google/generative-ai-go/genai"
)

func InitLLMService(cli *genai.Client) LLMService {
	return llm.NewGeminiService(cli)
}
