package discussion

import (
	"context"
	"fmt"

	"local-ai/internal/domain/discussion"
)

// GetDiscussionHistoryUseCase handles retrieving the message history of a discussion.
type GetDiscussionHistoryUseCase struct {
	repo discussion.Repository
}

// NewGetDiscussionHistoryUseCase creates a new get discussion history use case.
func NewGetDiscussionHistoryUseCase(repo discussion.Repository) *GetDiscussionHistoryUseCase {
	return &GetDiscussionHistoryUseCase{repo: repo}
}

// GetDiscussionHistoryQuery contains the parameters for getting discussion history.
type GetDiscussionHistoryQuery struct {
	DiscussionID string
}

// GetDiscussionHistoryResult contains the message history of a discussion.
type GetDiscussionHistoryResult struct {
	DiscussionID string       `json:"discussion_id"`
	Messages     []MessageDTO `json:"messages"`
}

// Execute retrieves the message history for a discussion.
func (uc *GetDiscussionHistoryUseCase) Execute(ctx context.Context, query GetDiscussionHistoryQuery) (*GetDiscussionHistoryResult, error) {
	if query.DiscussionID == "" {
		return nil, fmt.Errorf("discussion ID cannot be empty")
	}

	disc, err := uc.repo.FindByID(ctx, query.DiscussionID)
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

	return &GetDiscussionHistoryResult{
		DiscussionID: disc.ID,
		Messages:     messageDTOs,
	}, nil
}
