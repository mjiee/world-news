package openai

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// OpenaiClient is the client for the OpenAI API
type OpenaiClient struct {
	config *Config
	client *http.Client
}

// NewOpenaiClient creates a new OpenaiClient
func NewOpenaiClient(config *Config) *OpenaiClient {
	return &OpenaiClient{
		config: config,
		client: http.DefaultClient,
	}
}

// Do sends a request to the OpenAI API
func (c *OpenaiClient) Do(ctx context.Context, data []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, c.config.ApiUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.ApiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)

	return resp, errors.WithStack(err)
}
