package memory

import (
	"context"
	"fmt"
	"sync"

	"local-ai/internal/domain/discussion"
)

// InMemoryDiscussionRepository is an in-memory implementation of the discussion repository.
// This is suitable for development and testing. For production, replace with a database implementation.
type InMemoryDiscussionRepository struct {
	mu          sync.RWMutex
	discussions map[string]*discussion.Discussion
}

// NewInMemoryDiscussionRepository creates a new in-memory discussion repository.
func NewInMemoryDiscussionRepository() *InMemoryDiscussionRepository {
	return &InMemoryDiscussionRepository{
		discussions: make(map[string]*discussion.Discussion),
	}
}

// Save persists a discussion aggregate.
func (r *InMemoryDiscussionRepository) Save(ctx context.Context, disc *discussion.Discussion) error {
	if disc == nil {
		return fmt.Errorf("cannot save nil discussion")
	}
	if disc.ID == "" {
		return fmt.Errorf("cannot save discussion with empty ID")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Create a deep copy to prevent external mutations
	r.discussions[disc.ID] = r.copyDiscussion(disc)
	return nil
}

// FindByID retrieves a discussion by its unique identifier.
func (r *InMemoryDiscussionRepository) FindByID(ctx context.Context, id string) (*discussion.Discussion, error) {
	if id == "" {
		return nil, fmt.Errorf("discussion ID cannot be empty")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	disc, exists := r.discussions[id]
	if !exists {
		return nil, fmt.Errorf("discussion with ID %s not found", id)
	}

	// Return a copy to prevent external mutations
	return r.copyDiscussion(disc), nil
}

// FindAll retrieves all discussions.
func (r *InMemoryDiscussionRepository) FindAll(ctx context.Context) ([]*discussion.Discussion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	discussions := make([]*discussion.Discussion, 0, len(r.discussions))
	for _, disc := range r.discussions {
		discussions = append(discussions, r.copyDiscussion(disc))
	}

	return discussions, nil
}

// Delete removes a discussion by its ID.
func (r *InMemoryDiscussionRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("discussion ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.discussions[id]; !exists {
		return fmt.Errorf("discussion with ID %s not found", id)
	}

	delete(r.discussions, id)
	return nil
}

// Exists checks if a discussion with the given ID exists.
func (r *InMemoryDiscussionRepository) Exists(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("discussion ID cannot be empty")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.discussions[id]
	return exists, nil
}

// copyDiscussion creates a deep copy of a discussion to prevent mutations.
func (r *InMemoryDiscussionRepository) copyDiscussion(src *discussion.Discussion) *discussion.Discussion {
	if src == nil {
		return nil
	}

	// Copy messages
	messages := make([]*discussion.Message, len(src.Messages))
	for i, msg := range src.Messages {
		messages[i] = &discussion.Message{
			ID:        msg.ID,
			Role:      msg.Role,
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
		}
	}

	return &discussion.Discussion{
		ID:        src.ID,
		Title:     src.Title,
		Messages:  messages,
		Status:    src.Status,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}
