package documents

import "context"

type QueryResult struct {
	ID       string
	Distance float64
	Document string
	Metadata map[string]interface{}
}

type QueryResponse struct {
	Results []QueryResult
}

type DocumentRepository interface {
	AddDocument(ctx context.Context, doc *Document) (bool, error)

	RemoveDocument(ctx context.Context, docID string) (bool, error)

	Query(ctx context.Context, embedding []float64, limit int) (*QueryResponse, error)

	DocumentExists(ctx context.Context, docID string) (bool, error)
}
