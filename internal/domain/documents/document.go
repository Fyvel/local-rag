package documents

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"
)

type Document struct {
	ID         string
	Content    string
	SHA        string
	Embeddings []float64
	Metadata   map[string]interface{}
}

func NewDocument(id, content string, embeddings []float64, metadata map[string]interface{}) (*Document, error) {
	if id == "" {
		return nil, fmt.Errorf("document ID cannot be empty")
	}
	if content == "" {
		return nil, fmt.Errorf("document content cannot be empty")
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("document embeddings cannot be empty")
	}

	return &Document{
		ID:         id,
		Content:    content,
		Embeddings: embeddings,
		Metadata:   metadata,
	}, nil
}

func GenerateID() string {
	// Create buffer for 16 random bytes (similar to UUID)
	b := make([]byte, 16)

	// Read random bytes from crypto/rand
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("id-%d", time.Now().UnixNano())
	}

	// Format similar to UUID but without hyphens
	return fmt.Sprintf("%x", b)
}

// GenerateDeterministicID creates a deterministic identifier based on input string.
// This ensures the same input always generates the same ID, useful for deduplication.
func GenerateDeterministicID(input string) string {
	// Create a hash of the input string
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}
