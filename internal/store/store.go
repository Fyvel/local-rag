package store

import (
	"context"
	"fmt"

	"local-ai/internal/domain/documents"
	"local-ai/internal/store/chroma"
)

type VectorStore struct {
	client       *chroma.Client
	collectionID string
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

func (vs *VectorStore) AddDocument(ctx context.Context, doc *documents.Document) (bool, error) {
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
