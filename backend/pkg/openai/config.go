package openai

// Config is the configuration for the OpenAI API
type Config struct {
	Description string `json:"description"`

	ApiKey string `json:"apiKey"`
	ApiUrl string `json:"apiUrl"`

	Model     string `json:"model"`
	MaxTokens int    `json:"maxTokens"`

	SystemPrompt    string `json:"systemPrompt"`
	AssistantPrompt string `json:"assistantPrompt"`
}
