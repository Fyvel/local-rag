package embeddings

import "context"

// Embedder is a domain interface for services that generate embeddings.
// This belongs in the domain layer as it defines a core capability needed by business logic.
// Infrastructure implementations provide concrete embedding providers (e.g., Ollama).
type Embedder interface {
	// Embed generates embedding vectors for the given text.
	// Returns a slice of Embedding objects containing the vector representations.
	Embed(ctx context.Context, text string, model string) ([]*Embedding, error)
}
