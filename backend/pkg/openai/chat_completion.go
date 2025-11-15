package openai

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pkg/errors"

	"github.com/mjiee/world-news/backend/pkg/errorx"
)

// chatCompletionRequest is the request body for the OpenAI API
type chatCompletionRequest struct {
	Model     string     `json:"model"`
	Messages  []*Message `json:"messages"`
	MaxTokens int        `json:"max_tokens,omitempty"`
}

// SetSystemPrompt sets the system prompt for the ChatCompletionRequest
func (c *OpenaiClient) SetSystemPrompt(prompt string) *OpenaiClient {
	if c.chatCompletion == nil {
		c.chatCompletion = &chatCompletionRequest{}
	}

	if len(prompt) == 0 {
		return c
	}

	c.chatCompletion.Messages = append(c.chatCompletion.Messages, SystemMessage(prompt))

	return c
}

// SetUserPrompt sets the user prompt for the ChatCompletionRequest
func (c *OpenaiClient) SetUserPrompt(prompt ...string) *OpenaiClient {
	if c.chatCompletion == nil {
		c.chatCompletion = &chatCompletionRequest{}
	}

	if len(prompt) == 0 {
		return c
	}

	c.chatCompletion.Messages = append(c.chatCompletion.Messages, UserMessage(prompt...))

	return c
}

// SetAssistantPrompt sets the assistant prompt for the ChatCompletionRequest
func (c *OpenaiClient) SetAssistantPrompt(prompt string) *OpenaiClient {
	if c.chatCompletion == nil {
		c.chatCompletion = &chatCompletionRequest{}
	}

	if len(prompt) == 0 {
		return c
	}

	c.chatCompletion.Messages = append(c.chatCompletion.Messages, AssistantMessage(prompt))

	return c
}

// SetMessage sets the message for the ChatCompletionRequest
func (c *OpenaiClient) SetMessage(messages ...*Message) *OpenaiClient {
	if c.chatCompletion == nil {
		c.chatCompletion = &chatCompletionRequest{}
	}

	c.chatCompletion.Messages = append(c.chatCompletion.Messages, messages...)

	return c
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
func (c *OpenaiClient) ChatCompletion(ctx context.Context) (*ChatCompletionResponse, error) {
	if c.chatCompletion == nil {
		return nil, errors.WithStack(errorx.InternalError)
	}

	c.chatCompletion.Model = c.config.Model

	if c.config.MaxTokens > 0 {
		c.chatCompletion.MaxTokens = c.config.MaxTokens
	}

	reqBody, err := json.Marshal(c.chatCompletion)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := c.do(ctx, reqBody)
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
