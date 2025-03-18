package openai

import (
	"context"
	"encoding/json"
	"io"

	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/pkg/errors"
)

// ChatCompletionRequest is the request body for the OpenAI API
type ChatCompletionRequest struct {
	Model     string     `json:"model"`
	Messages  []*message `json:"messages"`
	MaxTokens int        `json:"max_tokens"`
}

// NewChatCompletionRequest creates a new ChatCompletionRequest
func NewChatCompletionRequest(conf *Config, userMsg []string) *ChatCompletionRequest {
	req := &ChatCompletionRequest{
		Model:     conf.Model,
		MaxTokens: conf.MaxTokens,
	}

	req.Messages = make([]*message, 0, 3)

	if conf.SystemPrompt != "" {
		req.Messages = append(req.Messages, SystemMessage(conf.SystemPrompt))
	}

	for _, msg := range userMsg {
		req.Messages = append(req.Messages, UserMessage(msg))
	}

	if conf.AssistantPrompt != "" {
		req.Messages = append(req.Messages, AssistantMessage(conf.AssistantPrompt))
	}

	return req
}

// ChatCompletionMessage is the message for the OpenAI API
type ChatCompletionMessage struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

// ChatCompletionResponse is the response body for the OpenAI API
type ChatCompletionResponse struct {
	ID                string                  `json:"id"`
	Model             string                  `json:"model"`
	Choices           []*ChatCompletionChoice `json:"choices"`
	Created           int                     `json:"created"`
	SystemFingerprint string                  `json:"system_fingerprint"`
	Object            string                  `json:"object"`
	Usage             *Usage                  `json:"usage,omitempty"`
	Error             *RespError              `json:"error,omitempty"`
}

// ChatCompletionChoice is the choice for the OpenAI API
type ChatCompletionChoice struct {
	Index        int                    `json:"index"`
	FinishReason string                 `json:"finish_reason"`
	Message      *ChatCompletionMessage `json:"message"`
}

// ChatCompletion is the chat completion endpoint for the OpenAI API
func (c *OpenaiClient) ChatCompletion(ctx context.Context, userMessage ...string) (*ChatCompletionResponse, error) {
	reqBody, err := json.Marshal(NewChatCompletionRequest(c.config, userMessage))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := c.Do(ctx, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result ChatCompletionResponse

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	if result.Error != nil && result.Error.Code != "" {
		return nil, errorx.InternalError.SetErr(errors.Errorf("%s: %s", result.Error.Code, result.Error.Message)).
			SetMessage(result.Error.Message)
	}

	return &result, nil
}
