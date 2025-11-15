package openai

import "strings"

// messageRole is the role for the OpenAI API
type messageRole string

const (
	System    messageRole = "system"
	User      messageRole = "user"
	Assistant messageRole = "assistant"
)

// Message is the Message for the OpenAI API
type Message struct {
	Role    messageRole `json:"role"`
	Content string      `json:"content"`
}

// UserMessage returns a new user message
func UserMessage(content ...string) *Message {
	return &Message{Role: User, Content: strings.Join(content, "\n")}
}

// SystemMessage returns a new system message
func SystemMessage(content string) *Message {
	return &Message{Role: System, Content: content}
}

// AssistantMessage returns a new assistant message
func AssistantMessage(content string) *Message {
	return &Message{Role: Assistant, Content: content}
}
