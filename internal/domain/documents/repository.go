package documents

import "context"

type DocumentRepository interface {
	AddDocument(ctx context.Context, doc *Document) (bool, error)

	RemoveDocument(ctx context.Context, docID string) (bool, error)

	Query(ctx context.Context, embedding []float64, limit int) (interface{}, error)
}
