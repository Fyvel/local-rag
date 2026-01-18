package indexing

import (
	"context"
	"fmt"
	"time"

	"local-ai/internal/domain/documents"
	"local-ai/internal/domain/embeddings"
	"local-ai/internal/domain/repositories"
	"local-ai/internal/domain/transformers"
)

// IndexRepositoryCommand contains the parameters for indexing a repository.
type IndexRepositoryCommand struct {
	GithubURL   string
	FileTypes   []string
	TargetIndex string
	Timeout     time.Duration
}

// IndexRepositoryResult contains the result of indexing a repository.
type IndexRepositoryResult struct {
	Repository *repositories.Repository `json:"repository"`
	IndexID    string                   `json:"index_id"`
}

// IndexRepositoryUseCase orchestrates the process of fetching and indexing a GitHub repository.
// This is an application service that coordinates between domain entities, repositories, and infrastructure.
type IndexRepositoryUseCase struct {
	repositoryFetcher repositories.Fetcher
	documentRepo      documents.DocumentRepository
	embedder          embeddings.Embedder
	textSplitter      transformers.TextChunker
	config            Config
}

// NewIndexRepositoryUseCase creates a new use case with injected dependencies.
// All dependencies are domain interfaces, making this truly infrastructure-agnostic.
func NewIndexRepositoryUseCase(
	repositoryFetcher repositories.Fetcher,
	documentRepo documents.DocumentRepository,
	embedder embeddings.Embedder,
	textSplitter transformers.TextChunker,
	config Config,
) *IndexRepositoryUseCase {
	return &IndexRepositoryUseCase{
		repositoryFetcher: repositoryFetcher,
		documentRepo:      documentRepo,
		embedder:          embedder,
		textSplitter:      textSplitter,
		config:            config,
	}
}

// Execute runs the index repository use case.
// It orchestrates fetching the repository and indexing its contents.
func (uc *IndexRepositoryUseCase) Execute(ctx context.Context, cmd IndexRepositoryCommand) (*IndexRepositoryResult, error) {
	if cmd.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cmd.Timeout)
		defer cancel()
	} else {
		// Use default timeout from config
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, uc.config.OperationTimeout)
		defer cancel()
	}

	fmt.Printf("Fetching repository: %s\n", cmd.GithubURL)
	repo, err := uc.repositoryFetcher.FetchRepository(ctx, cmd.GithubURL, cmd.FileTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository: %w", err)
	}

	fmt.Printf("Repository SHA: %s\n", repo.SHA)
	fmt.Printf("Files found: %d\n", repo.FileCount())

	idx := &Indexer{
		TextSplitter:    uc.textSplitter,
		VectorStore:     uc.documentRepo,
		Embedding:       uc.embedder,
		DocumentService: &documents.DocumentService{},
		Config:          uc.config,
	}

	fmt.Println("Starting indexing process...")
	processedRepo, err := idx.ProcessRepository(ctx, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to index repository: %w", err)
	}

	return &IndexRepositoryResult{
		Repository: processedRepo,
		IndexID:    cmd.TargetIndex,
	}, nil
}
