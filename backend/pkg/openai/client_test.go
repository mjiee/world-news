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

	response, err := client.ChatCompletion(context.Background(),
		"Hawks vs. Hornets Predictions: Odds, expert picks, recent stats, trends and best bets for March 18",
		"It’s Tuesday, March 18, and the Atlanta Hawks (32-36) and Charlotte Hornets (17-50) are all set to square off from Spectrum Center in Charlotte.\nThe Hawks are currently 15-18 on the road with a point differential of -3, while the Hornets have a 2-8 record in their last ten games at home. Atlanta is 3-0 this season against Charlotte with wins of 3, 5, and 13 points.",
	)
	if err != nil {
		t.Fatalf("ChatCompletion failed: %v", err)
	}

	for _, choice := range response.Choices {
		fmt.Println(choice.Message.Content)
	}
}
