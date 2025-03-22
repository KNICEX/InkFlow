package ai

import (
	"github.com/KNICEX/InkFlow/internal/ai/internal/service/llm/gemini"
	"github.com/google/generative-ai-go/genai"
)

func InitLLMService(cli *genai.Client) LLMService {
	return gemini.NewGeminiService(cli)
}
