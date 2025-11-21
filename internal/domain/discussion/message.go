package discussion

import (
	"errors"
	"time"
)

// MessageRole represents the role of the message sender.
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

// Message represents a single message in a discussion.
// It's a value object - immutable after creation.
type Message struct {
	ID        string
	Role      MessageRole
	Content   string
	Timestamp time.Time
}

// NewMessage creates a new message with validation.
func NewMessage(id string, role MessageRole, content string) (*Message, error) {
	if id == "" {
		return nil, errors.New("message ID cannot be empty")
	}
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}
	if role != RoleUser && role != RoleAssistant {
		return nil, errors.New("invalid message role")
	}

	return &Message{
		ID:        id,
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}, nil
}
