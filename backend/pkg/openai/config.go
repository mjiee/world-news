package openai

// Config is the configuration for the OpenAI API
type Config struct {
	Platform string `json:"platform"`

	ApiKey string `json:"apiKey"`
	ApiUrl string `json:"apiUrl"`

	Model     string `json:"model"`
	MaxTokens int    `json:"maxTokens"`
}
