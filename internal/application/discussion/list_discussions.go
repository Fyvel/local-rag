package discussion

import (
	"context"
	"fmt"

	"local-ai/internal/domain/discussion"
)

// ListDiscussionsUseCase handles listing all discussions.
type ListDiscussionsUseCase struct {
	repo discussion.Repository
}

// NewListDiscussionsUseCase creates a new list discussions use case.
func NewListDiscussionsUseCase(repo discussion.Repository) *ListDiscussionsUseCase {
	return &ListDiscussionsUseCase{repo: repo}
}

// Execute retrieves all discussions.
func (uc *ListDiscussionsUseCase) Execute(ctx context.Context) ([]DiscussionSummaryResult, error) {
	discussions, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve discussions: %w", err)
	}

	summaries := make([]DiscussionSummaryResult, 0, len(discussions))
	for _, d := range discussions {
		summaries = append(summaries, DiscussionSummaryResult{
			ID:           d.ID,
			Title:        d.Title,
			MessageCount: d.MessageCount(),
			Status:       string(d.Status),
			CreatedAt:    d.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    d.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return summaries, nil
}
