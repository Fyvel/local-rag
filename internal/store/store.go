package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"local-ai/internal/store/chroma"
)

type VectorStore struct {
	client       *chroma.Client
	collectionID string
}

type Document struct {
	ID         string
	Content    string
	SHA        string
	Embeddings []float64
	Metadata   map[string]interface{}
}

func NewVectorStore(opts ...chroma.Option) *VectorStore {
	client := chroma.NewClient(opts...)

	ctx := context.Background()
	collection, err := client.EnsureCollection(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to ensure collection: %v", err))
	}

	return &VectorStore{
		client:       client,
		collectionID: collection.ID,
	}
}

func (vs *VectorStore) AddDocument(ctx context.Context, doc Document) (bool, error) {
	req := &chroma.AddDocumentsRequest{
		IDs:        []string{doc.ID},
		Embeddings: [][]float64{doc.Embeddings},
		Metadatas:  []map[string]interface{}{doc.Metadata},
		Documents:  []string{doc.Content},
	}

	errDelete := vs.client.DeleteDocuments(ctx, vs.collectionID, []string{doc.ID})
	if errDelete != nil {
		return false, fmt.Errorf("failed to clean old document: %w", errDelete)
	}

	err := vs.client.AddDocuments(ctx, vs.collectionID, req)
	if err != nil {
		return false, fmt.Errorf("failed to add document: %w", err)
	}

	return true, nil
}

func (vs *VectorStore) RemoveDocument(docID string) (bool, error) {
	ctx := context.Background()
	err := vs.client.DeleteDocuments(ctx, vs.collectionID, []string{docID})
	if err != nil {
		return false, fmt.Errorf("failed to remove document: %w", err)
	}

	return true, nil
}

func (vs *VectorStore) Query(ctx context.Context, embedding []float64, limit int) ([]chroma.QueryResult, error) {
	return vs.client.Query(ctx, vs.collectionID, embedding, limit)
}

type RepositoryRecord struct {
	IndexId   string
	Url       string
	SHA       string
	UpdatedAt string
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

func GenerateDeterministicID(input string) string {
	// Create a hash of the input string
	hash := fmt.Sprintf("%x", input)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(hash)))
}
