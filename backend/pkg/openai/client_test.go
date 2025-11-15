package openai

import (
	"context"
	"fmt"
	"testing"
)

var (
	testConfig = Config{
		Platform: "test",
		ApiKey:   "xxx",
		ApiUrl:   "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
		Model:    "deepseek-r1",
	}
)

// TestChatCompletion tests the chat completion endpoint
func TestChatCompletion(t *testing.T) {
	response, err := NewOpenaiClient(&testConfig).
		SetSystemPrompt("You are a professional news commentary assistant.").
		SetUserPrompt("Hawks vs. Hornets Predictions: Odds, expert picks, recent stats, trends and best bets for March 18").
		ChatCompletion(context.Background())
	if err != nil {
		t.Fatalf("ChatCompletion failed: %v", err)
	}

	for _, choice := range response.Choices {
		fmt.Println(choice.Message.Content)
	}
}
