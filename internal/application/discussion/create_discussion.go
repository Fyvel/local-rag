package discussion

import (
	"context"
	"fmt"

	"local-ai/internal/domain/discussion"

	"github.com/google/uuid"
)

// CreateDiscussionUseCase handles creating a new discussion.
type CreateDiscussionUseCase struct {
	repo discussion.Repository
}

// NewCreateDiscussionUseCase creates a new create discussion use case.
func NewCreateDiscussionUseCase(repo discussion.Repository) *CreateDiscussionUseCase {
	return &CreateDiscussionUseCase{repo: repo}
}

// CreateDiscussionCommand contains the parameters for creating a discussion.
type CreateDiscussionCommand struct {
	Title string
}

// CreateDiscussionResult contains the result of creating a discussion.
type CreateDiscussionResult struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// Execute creates a new discussion.
func (uc *CreateDiscussionUseCase) Execute(ctx context.Context, cmd CreateDiscussionCommand) (*CreateDiscussionResult, error) {
	if cmd.Title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	// Generate a unique ID for the discussion
	id := uuid.New().String()

	// Create the discussion aggregate
	disc, err := discussion.NewDiscussion(id, cmd.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to create discussion: %w", err)
	}

	// Persist the discussion
	if err := uc.repo.Save(ctx, disc); err != nil {
		return nil, fmt.Errorf("failed to save discussion: %w", err)
	}

	return &CreateDiscussionResult{
		ID:        disc.ID,
		Title:     disc.Title,
		Status:    string(disc.Status),
		CreatedAt: disc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
