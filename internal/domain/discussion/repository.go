package discussion

import "context"

// Repository defines the interface for discussion persistence.
// This is a domain interface that will be implemented by the infrastructure layer.
type Repository interface {
	// Save persists a discussion aggregate.
	Save(ctx context.Context, discussion *Discussion) error

	// FindByID retrieves a discussion by its unique identifier.
	FindByID(ctx context.Context, id string) (*Discussion, error)

	// FindAll retrieves all discussions.
	// Returns a slice of discussions, potentially empty if none exist.
	FindAll(ctx context.Context) ([]*Discussion, error)

	// Delete removes a discussion by its ID.
	Delete(ctx context.Context, id string) error

	// Exists checks if a discussion with the given ID exists.
	Exists(ctx context.Context, id string) (bool, error)
}
