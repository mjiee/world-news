package openai

import (
	"context"
	"fmt"
	"testing"
)

var (
	testConfig = Config{
		Description:  "test",
		ApiKey:       "xxx",
		ApiUrl:       "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
		Model:        "deepseek-r1",
		SystemPrompt: "You are a professional news commentary assistant.",
	}
)

// TestChatCompletion tests the chat completion endpoint
func TestChatCompletion(t *testing.T) {
	client := NewOpenaiClient(&testConfig)

	response, err := client.ChatCompletion(context.Background(), "Duty on steel and aluminium imports is a major hit for some of the US's top trading partners.")
	if err != nil {
		t.Fatalf("ChatCompletion failed: %v", err)
	}

	fmt.Println(response)
}
