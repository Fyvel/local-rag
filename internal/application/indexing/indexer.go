package indexing

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"local-ai/internal/domain/documents"
	"local-ai/internal/domain/embeddings"
	"local-ai/internal/domain/repositories"
	"local-ai/internal/domain/transformers"
)

type Config struct {
	WorkerCount      int
	ChunkMaxSize     int
	ChunkOverlap     int
	OperationTimeout time.Duration
	EmbeddingModel   string
}

// DefaultConfig provides sensible defaults
func DefaultConfig() Config {
	return Config{
		WorkerCount:      runtime.NumCPU() / 2,
		ChunkMaxSize:     1000,
		ChunkOverlap:     100,
		OperationTimeout: 10 * time.Minute,
		EmbeddingModel:   "mxbai-embed-large",
	}
}

type Indexer struct {
	TextSplitter    transformers.TextChunker
	VectorStore     documents.DocumentRepository
	Embedding       embeddings.Embedder
	DocumentService *documents.DocumentService
	Config          Config
}

func (indexer *Indexer) ProcessRepository(ctx context.Context, repository *repositories.Repository) (*repositories.Repository, error) {
	var wg sync.WaitGroup
	filesChan := make(chan repositories.File, len(repository.Files))
	errChan := make(chan error, len(repository.Files))

	// Prepare error tracking
	var (
		mu       sync.Mutex
		firstErr error
	)

	// Start workers
	for range indexer.Config.WorkerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range filesChan {
				// Check if context is done
				select {
				case <-ctx.Done():
					mu.Lock()
					if firstErr == nil {
						firstErr = ctx.Err()
					}
					mu.Unlock()
					return
				default:
					if err := indexer.processFile(ctx, repository, file); err != nil {
						mu.Lock()
						if firstErr == nil {
							firstErr = err
						}
						errChan <- err
						mu.Unlock()
					}
				}
			}
		}()
	}

	// Send files to workers
	for _, file := range repository.Files {
		filesChan <- file
	}
	close(filesChan)

	// Wait for all workers to finish
	wg.Wait()
	close(errChan)

	// Collect all errors
	var allErrors []error
	for err := range errChan {
		allErrors = append(allErrors, err)
	}

	// Report total errors
	if len(allErrors) > 0 {
		fmt.Printf("Encountered %d errors during processing\n", len(allErrors))
		return repository, fmt.Errorf("errors processing repository: %v", allErrors)
	}

	fmt.Println("Processing complete")

	return repository, nil
}

// processFile handles indexing of an individual file
func (indexer *Indexer) processFile(ctx context.Context, repository *repositories.Repository, file repositories.File) error {
	// Context check
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		fmt.Printf("Processing file: %s\n", file.Name)
		fmt.Printf("File path: %s\n", file.Path)

		chunks, err := indexer.TextSplitter.SplitText(file.Content)
		if err != nil {
			return fmt.Errorf("failed to split file %s: %w", file.Name, err)
		}

		for chunkIndex, chunk := range chunks {
			// Generate embedding for the chunk
			embedding, err := indexer.Embedding.Embed(ctx, chunk, indexer.Config.EmbeddingModel)
			if err != nil {
				return fmt.Errorf("embedding failed for file %s, chunk %d: %w", file.Name, chunkIndex, err)
			}

			// Use domain service to create document from chunk
			metadata := documents.ChunkMetadata{
				FileName:       file.Name,
				FilePath:       file.Path,
				FileType:       file.Type,
				FileSize:       len(file.Content),
				ChunkIndex:     chunkIndex,
				TotalChunks:    len(chunks),
				RepositoryURL:  repository.URL,
				RepositoryName: repository.Name,
				CommitSHA:      repository.SHA,
			}

			doc, err := indexer.DocumentService.CreateDocumentFromChunk(
				chunk,
				embedding[0].Vector,
				metadata,
			)
			if err != nil {
				return fmt.Errorf("failed to create document for file %s, chunk %d: %w", file.Name, chunkIndex, err)
			}

			// Validate document using domain service
			if err := indexer.DocumentService.ValidateDocument(doc); err != nil {
				return fmt.Errorf("document validation failed for file %s, chunk %d: %w", file.Name, chunkIndex, err)
			}

			// Store document
			if _, err := indexer.VectorStore.AddDocument(ctx, doc); err != nil {
				return fmt.Errorf("failed to add document to vector store: %w", err)
			}
		}

		return nil
	}
}
