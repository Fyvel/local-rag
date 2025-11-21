package discussion

import (
	"context"
	"fmt"

	"local-ai/internal/domain/discussion"
)

// GetDiscussionUseCase handles retrieving a discussion by ID.
type GetDiscussionUseCase struct {
	repo discussion.Repository
}

// NewGetDiscussionUseCase creates a new get discussion use case.
func NewGetDiscussionUseCase(repo discussion.Repository) *GetDiscussionUseCase {
	return &GetDiscussionUseCase{repo: repo}
}

// MessageDTO is a data transfer object for messages.
type MessageDTO struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// DiscussionDTO is a data transfer object for a complete discussion.
type DiscussionDTO struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Messages  []MessageDTO `json:"messages"`
	Status    string       `json:"status"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}

// GetDiscussionQuery contains the parameters for getting a discussion.
type GetDiscussionQuery struct {
	ID string
}

// Execute retrieves a discussion by ID.
func (uc *GetDiscussionUseCase) Execute(ctx context.Context, query GetDiscussionQuery) (*DiscussionDTO, error) {
	if query.ID == "" {
		return nil, fmt.Errorf("discussion ID cannot be empty")
	}

	disc, err := uc.repo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve discussion: %w", err)
	}

	// Convert messages to DTOs
	messageDTOs := make([]MessageDTO, 0, len(disc.Messages))
	for _, msg := range disc.Messages {
		messageDTOs = append(messageDTOs, MessageDTO{
			ID:        msg.ID,
			Role:      string(msg.Role),
			Content:   msg.Content,
			Timestamp: msg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &DiscussionDTO{
		ID:        disc.ID,
		Title:     disc.Title,
		Messages:  messageDTOs,
		Status:    string(disc.Status),
		CreatedAt: disc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: disc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
