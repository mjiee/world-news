package openai

// messageRole is the role for the OpenAI API
type messageRole string

const (
	system    messageRole = "system"
	user      messageRole = "user"
	assistant messageRole = "assistant"
)

// message is the message for the OpenAI API
type message struct {
	Role    messageRole `json:"role"`
	Content string      `json:"content"`
}

// UserMessage returns a new user message
func UserMessage(content string) *message {
	return &message{Role: user, Content: content}
}

// SystemMessage returns a new system message
func SystemMessage(content string) *message {
	return &message{Role: system, Content: content}
}

// AssistantMessage returns a new assistant message
func AssistantMessage(content string) *message {
	return &message{Role: assistant, Content: content}
}
