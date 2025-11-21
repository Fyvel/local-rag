package discussion

import (
	"errors"
	"time"
)

// DiscussionStatus represents the current state of a discussion.
type DiscussionStatus string

const (
	StatusOpen   DiscussionStatus = "open"
	StatusClosed DiscussionStatus = "closed"
)

// Discussion is the aggregate root for discussion domain.
// It maintains the list of messages and enforces domain invariants.
type Discussion struct {
	ID        string
	Title     string
	Messages  []*Message
	Status    DiscussionStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewDiscussion creates a new discussion aggregate.
func NewDiscussion(id, title string) (*Discussion, error) {
	if id == "" {
		return nil, errors.New("discussion ID cannot be empty")
	}
	if title == "" {
		return nil, errors.New("discussion title cannot be empty")
	}

	now := time.Now()
	return &Discussion{
		ID:        id,
		Title:     title,
		Messages:  make([]*Message, 0),
		Status:    StatusOpen,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// AddMessage adds a message to the discussion.
// Enforces domain invariants: no empty messages, no messages on closed discussions.
func (d *Discussion) AddMessage(message *Message) error {
	if message == nil {
		return errors.New("cannot add nil message")
	}
	if message.Content == "" {
		return errors.New("cannot add empty message")
	}
	if d.Status == StatusClosed {
		return errors.New("cannot add message to closed discussion")
	}

	d.Messages = append(d.Messages, message)
	d.UpdatedAt = time.Now()
	return nil
}

// Close closes the discussion, preventing further messages.
func (d *Discussion) Close() error {
	if d.Status == StatusClosed {
		return errors.New("discussion is already closed")
	}

	d.Status = StatusClosed
	d.UpdatedAt = time.Now()
	return nil
}

// GetMessages returns a copy of the messages to prevent external mutation.
func (d *Discussion) GetMessages() []*Message {
	messages := make([]*Message, len(d.Messages))
	copy(messages, d.Messages)
	return messages
}

// MessageCount returns the number of messages in the discussion.
func (d *Discussion) MessageCount() int {
	return len(d.Messages)
}
