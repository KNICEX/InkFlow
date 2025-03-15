package service

import (
	"context"
	"fmt"
	"testing"
)

func TestGeminiService_AskOnce(t *testing.T) {
	svc, err := NewGeminiService("AIzaSyDsK-uD5Y-mW17slUROmaA4kFpohM4V96Y")
	if err != nil {
		panic(err)
	}
	ask, err := svc.AskOnce(context.Background(), "你是谁？")
	if err != nil {
		panic(err)
	}
	fmt.Println(ask)
}
