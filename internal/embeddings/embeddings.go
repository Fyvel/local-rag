package embeddings

import (
	"context"

	"local-ai/internal/domain/embeddings"
)

// Embedder is an interface for services that generate embeddings.
// This is an infrastructure-level interface for embedding providers.
type Embedder[T any] interface {
	Embed(context.Context, T) ([]*embeddings.Embedding, error)
}
