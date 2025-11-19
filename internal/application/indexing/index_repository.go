package indexing

import (
	"context"
	"fmt"
	"time"

	"github.com/tmc/langchaingo/textsplitter"

	"local-ai/internal/domain/documents"
	"local-ai/internal/domain/repositories"
	"local-ai/internal/infra/embeddings"
	"local-ai/internal/infra/embeddings/ollama"
	"local-ai/internal/infra/fetcher"
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
	Repository *repositories.Repository
	IndexID    string
}

// IndexRepositoryUseCase orchestrates the process of fetching and indexing a GitHub repository.
// This is an application service that coordinates between domain entities, repositories, and infrastructure.
type IndexRepositoryUseCase struct {
	repositoryFetcher RepositoryFetcher
	documentRepo      documents.DocumentRepository
	embedder          embeddings.Embedder[*ollama.EmbeddingRequest]
	textSplitter      *textsplitter.MarkdownTextSplitter
	config            Config
}

// RepositoryFetcher defines the contract for fetching repositories.
// This allows the use case to be independent of the specific fetcher implementation.
type RepositoryFetcher interface {
	GetGithubRepository(ctx context.Context, githubURL string, fileTypes []string) (*repositories.Repository, error)
}

// NewIndexRepositoryUseCase creates a new use case with injected dependencies.
func NewIndexRepositoryUseCase(
	documentRepo documents.DocumentRepository,
	embedder embeddings.Embedder[*ollama.EmbeddingRequest],
	textSplitter *textsplitter.MarkdownTextSplitter,
	config Config,
) *IndexRepositoryUseCase {
	return &IndexRepositoryUseCase{
		repositoryFetcher: &gitHubFetcher{},
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
	repo, err := uc.repositoryFetcher.GetGithubRepository(ctx, cmd.GithubURL, cmd.FileTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository: %w", err)
	}

	fmt.Printf("Repository SHA: %s\n", repo.SHA)
	fmt.Printf("Files found: %d\n", repo.FileCount())

	idx := &Indexer{
		TextSplitter:    uc.textSplitter,
		VectorStore:     uc.documentRepo,
		Embedding:       uc.embedder,
		DocumentService: documents.NewDocumentService(),
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

type gitHubFetcher struct{}

func (f *gitHubFetcher) GetGithubRepository(ctx context.Context, githubURL string, fileTypes []string) (*repositories.Repository, error) {
	return fetcher.GetGithubRepository(ctx, githubURL, fileTypes)
}
