package vectorstore

import (
	"context"
	"fmt"

	"local-ai/internal/domain/documents"
	"local-ai/internal/infra/store/chroma"
)

// VectorStore is an infrastructure implementation of the DocumentRepository interface.
// It uses ChromaDB as the underlying vector database.
type VectorStore struct {
	client       *chroma.Client
	collectionID string
}

var _ documents.DocumentRepository = (*VectorStore)(nil)

func New(opts ...chroma.Option) *VectorStore {
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

func (vs *VectorStore) RemoveDocument(ctx context.Context, docID string) (bool, error) {
	err := vs.client.DeleteDocuments(ctx, vs.collectionID, []string{docID})
	if err != nil {
		return false, fmt.Errorf("failed to remove document: %w", err)
	}

	return true, nil
}

func (vs *VectorStore) Query(ctx context.Context, embedding []float64, limit int) (*documents.QueryResponse, error) {
	results, err := vs.client.Query(ctx, vs.collectionID, embedding, limit)
	if err != nil {
		return nil, err
	}

	// Convert Chroma results to domain QueryResults
	domainResults := make([]documents.QueryResult, len(results))
	for i, result := range results {
		domainResults[i] = documents.QueryResult{
			ID:       result.ID,
			Distance: result.Distance,
			Document: result.Document,
			Metadata: result.Metadata,
		}
	}

	return &documents.QueryResponse{
		Results: domainResults,
	}, nil
}
