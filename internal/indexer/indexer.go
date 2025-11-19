package indexer

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/tmc/langchaingo/textsplitter"

	"local-ai/internal/client"
	"local-ai/internal/domain/documents"
	"local-ai/internal/domain/repositories"
	"local-ai/internal/embeddings"
	"local-ai/internal/embeddings/ollama"
	"local-ai/internal/fetcher"
	"local-ai/internal/store"
	"local-ai/internal/store/chroma"
)

type IndexerConfig struct {
	WorkerCount      int
	ChunkMaxSize     int
	ChunkOverlap     int
	OperationTimeout time.Duration
}

// DefaultConfig provides sensible defaults
func DefaultConfig() IndexerConfig {
	return IndexerConfig{
		WorkerCount:      runtime.NumCPU() / 2,
		ChunkMaxSize:     1000,
		ChunkOverlap:     100,
		OperationTimeout: 10 * time.Minute,
	}
}

type Indexer struct {
	TextSplitter *textsplitter.MarkdownTextSplitter
	VectorStore  *store.VectorStore
	Embedding    embeddings.Embedder[*ollama.EmbeddingRequest]
	Config       IndexerConfig
}

func IndexGithubRepository(
	ctx context.Context,
	githubURL string,
	fileTypes []string,
	targetIndex string,
) (*repositories.Repository, error) {
	config := DefaultConfig()

	ctx, cancel := context.WithTimeout(ctx, config.OperationTimeout)
	defer cancel()

	// Fetch GitHub repository
	repo, err := fetcher.GetGithubRepository(ctx, githubURL, fileTypes)
	if err != nil {
		return repo, fmt.Errorf("failed to fetch repository: %w", err)
	}

	fmt.Printf("Repository SHA: %s\n", repo.SHA)
	fmt.Printf("Files found: %d\n", len(repo.Files))

	fmt.Printf("	>	Init text splitter...\n")
	splitter := textsplitter.NewMarkdownTextSplitter([]textsplitter.Option{
		// textsplitter.WithMaxCharacters(cfg.ChunkMaxSize),
		// textsplitter.WithOverlap(cfg.ChunkOverlap),
	}...)

	fmt.Print("	>	Init vector store...\n")
	vs := store.NewVectorStore(
		chroma.WithChromaURL("http://localhost:8000"),
		chroma.WithHTTPClient(client.NewHTTP()),
		chroma.WithCollection(targetIndex),
	)

	fmt.Print("	>	Init embedding...\n")
	embedding := ollama.NewEmbedder(
		ollama.WithBaseURL("http://localhost:11434/api"),
		ollama.WithHTTPClient(client.NewHTTP()),
	)

	indexer := &Indexer{
		TextSplitter: splitter,
		VectorStore:  vs,
		Embedding:    embedding,
		Config:       config,
		// Logger:       logger,
	}

	return indexer.processRepository(ctx, repo)
}

func (indexer *Indexer) processRepository(ctx context.Context, repository *repositories.Repository) (*repositories.Repository, error) {
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
			embedding, err := indexer.Embedding.Embed(ctx, &ollama.EmbeddingRequest{
				Prompt: chunk,
				Model:  "mxbai-embed-large",
			})

			if err != nil {
				return fmt.Errorf("embedding failed for file %s, chunk %d: %w", file.Name, chunkIndex, err)
			}

			docIdStr := fmt.Sprintf("%s:%s:%s:%d", repository.Name, repository.SHA, file.Path, chunkIndex)

			doc := &documents.Document{
				ID:         documents.GenerateDeterministicID(docIdStr),
				Content:    chunk,
				Embeddings: embedding[0].Vector,
				Metadata: map[string]any{
					"file":        file.Name,
					"filePath":    file.Path,
					"fileType":    file.Type,
					"fileSize":    len(file.Content),
					"chunkIndex":  chunkIndex,
					"totalChunks": len(chunks),
					"githubURL":   repository.URL,
					"repoName":    repository.Name,
					"commitSHA":   repository.SHA,
					"indexedAt":   time.Now().Format(time.RFC3339),
				},
			}

			if _, err := indexer.VectorStore.AddDocument(ctx, doc); err != nil {
				return fmt.Errorf("failed to add document to vector store: %w", err)
			}
		}

		return nil
	}
}
