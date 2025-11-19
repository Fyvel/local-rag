package indexing

import (
	"context"
	"fmt"
	"time"

	"github.com/tmc/langchaingo/textsplitter"

	"local-ai/internal/client"
	"local-ai/internal/domain/documents"
	"local-ai/internal/domain/repositories"
	"local-ai/internal/embeddings"
	"local-ai/internal/embeddings/ollama"
	"local-ai/internal/fetcher"
	"local-ai/internal/indexer"
	"local-ai/internal/infra/store/chroma"
	"local-ai/internal/store"
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
	config            indexer.IndexerConfig
}

// RepositoryFetcher defines the contract for fetching repositories.
// This allows the use case to be independent of the specific fetcher implementation.
type RepositoryFetcher interface {
	GetGithubRepository(ctx context.Context, githubURL string, fileTypes []string) (*repositories.Repository, error)
}

// NewIndexRepositoryUseCase creates a new use case with default dependencies.
// In a production application, these dependencies would be injected via DI container.
func NewIndexRepositoryUseCase(
	targetIndex string,
	chromaURL string,
	ollamaURL string,
) *IndexRepositoryUseCase {
	// Initialize text splitter
	textSplitter := textsplitter.NewMarkdownTextSplitter([]textsplitter.Option{
		// Can be configured with options
	}...)

	// Initialize vector store
	vectorStore := store.NewVectorStore(
		chroma.WithChromaURL(chromaURL),
		chroma.WithHTTPClient(client.NewHTTP()),
		chroma.WithCollection(targetIndex),
	)

	// Initialize embedder
	embedder := ollama.NewEmbedder(
		ollama.WithBaseURL(ollamaURL),
		ollama.WithHTTPClient(client.NewHTTP()),
	)

	return &IndexRepositoryUseCase{
		repositoryFetcher: &gitHubFetcher{},
		documentRepo:      vectorStore,
		embedder:          embedder,
		textSplitter:      textSplitter,
		config:            indexer.DefaultConfig(),
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

	idx := &indexer.Indexer{
		TextSplitter: uc.textSplitter,
		VectorStore:  uc.documentRepo,
		Embedding:    uc.embedder,
		Config:       uc.config,
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
