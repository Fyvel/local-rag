package documents

import (
	"fmt"
	"time"
)

// ChunkMetadata contains metadata about a document chunk's origin.
type ChunkMetadata struct {
	FileName       string
	FilePath       string
	FileType       string
	FileSize       int
	ChunkIndex     int
	TotalChunks    int
	RepositoryURL  string
	RepositoryName string
	CommitSHA      string
}

// DocumentService provides domain logic for document operations.
// This is a domain service that encapsulates business rules for document creation and management.
type DocumentService struct{}

// NewDocumentService creates a new DocumentService.
func NewDocumentService() *DocumentService {
	return &DocumentService{}
}

// CreateDocumentFromChunk creates a Document entity from a text chunk with its metadata.
// This encapsulates the business logic of how documents are constructed from chunks,
// including ID generation, metadata enrichment, and timestamping.
func (s *DocumentService) CreateDocumentFromChunk(
	chunk string,
	embedding []float64,
	metadata ChunkMetadata,
) (*Document, error) {
	// Validate inputs
	if chunk == "" {
		return nil, fmt.Errorf("chunk content cannot be empty")
	}
	if len(embedding) == 0 {
		return nil, fmt.Errorf("embedding vector cannot be empty")
	}

	// Generate deterministic ID based on repository, file, and chunk position
	// This ensures the same chunk always gets the same ID for idempotency
	docIDStr := fmt.Sprintf("%s:%s:%s:%d",
		metadata.RepositoryName,
		metadata.CommitSHA,
		metadata.FilePath,
		metadata.ChunkIndex,
	)
	docID := GenerateDeterministicID(docIDStr)

	// Build comprehensive metadata
	docMetadata := map[string]interface{}{
		"file":        metadata.FileName,
		"filePath":    metadata.FilePath,
		"fileType":    metadata.FileType,
		"fileSize":    metadata.FileSize,
		"chunkIndex":  metadata.ChunkIndex,
		"totalChunks": metadata.TotalChunks,
		"githubURL":   metadata.RepositoryURL,
		"repoName":    metadata.RepositoryName,
		"commitSHA":   metadata.CommitSHA,
		"indexedAt":   time.Now().Format(time.RFC3339),
	}

	// Create document entity
	return &Document{
		ID:         docID,
		Content:    chunk,
		Embeddings: embedding,
		Metadata:   docMetadata,
	}, nil
}

// ValidateDocument checks if a document meets domain requirements.
// This encapsulates business rules for what constitutes a valid document.
func (s *DocumentService) ValidateDocument(doc *Document) error {
	if doc == nil {
		return fmt.Errorf("document cannot be nil")
	}
	if doc.ID == "" {
		return fmt.Errorf("document ID cannot be empty")
	}
	if doc.Content == "" {
		return fmt.Errorf("document content cannot be empty")
	}
	if len(doc.Embeddings) == 0 {
		return fmt.Errorf("document embeddings cannot be empty")
	}
	return nil
}
