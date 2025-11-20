package documents

import (
	"crypto/sha256"
	"fmt"
)

type Document struct {
	ID         string
	Content    string
	Embeddings []float64
	Metadata   map[string]interface{}
}

// GenerateDeterministicID creates a deterministic identifier based on input string.
// This ensures the same input always generates the same ID, useful for deduplication.
func GenerateDeterministicID(input string) string {
	// Create a hash of the input string
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}
