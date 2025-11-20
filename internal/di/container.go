package di

import (
	"local-ai/internal/application/indexing"
	"local-ai/internal/domain/documents"
	"local-ai/internal/domain/embeddings"
	"local-ai/internal/domain/repositories"
	"local-ai/internal/domain/transformers"
	"local-ai/internal/infra/embeddings/ollama"
	"local-ai/internal/infra/fetcher"
	"local-ai/internal/infra/httpclient"
	"local-ai/internal/infra/store/chroma"
	"local-ai/internal/infra/store/vectorstore"
	"local-ai/internal/infra/textsplitter"
)

// Config holds the configuration for dependency injection.
type Config struct {
	ChromaURL    string
	OllamaURL    string
	TargetIndex  string
	ChunkMaxSize int
	ChunkOverlap int
}

// Container holds all application dependencies.
type Container struct {
	HTTPClient        *httpclient.Client
	VectorStore       documents.DocumentRepository
	Embedder          embeddings.Embedder
	RepositoryFetcher repositories.Fetcher
	TextSplitter      transformers.TextChunker
	DocumentService   *documents.DocumentService
	IndexingConfig    indexing.Config
}

// NewContainer creates and wires up all application dependencies.
func NewContainer(cfg Config) *Container {
	// Initialize HTTP client (shared across services)
	httpClient := httpclient.New()

	// Initialize text splitter (infrastructure adapter for domain interface)
	textSplitter := textsplitter.NewMarkdownSplitter()

	// Initialize vector store
	vectorStore := vectorstore.New(
		chroma.WithChromaURL(cfg.ChromaURL),
		chroma.WithHTTPClient(httpClient),
		chroma.WithCollection(cfg.TargetIndex),
	)

	// Initialize embedder
	embedder := ollama.NewEmbedder(
		ollama.WithBaseURL(cfg.OllamaURL),
		ollama.WithHTTPClient(httpClient),
	)

	// Initialize repository fetcher
	repositoryFetcher := fetcher.NewGitHubFetcher()

	// Initialize domain services
	documentService := documents.NewDocumentService()

	// Initialize indexing config
	indexingConfig := indexing.DefaultConfig()

	return &Container{
		HTTPClient:        httpClient,
		VectorStore:       vectorStore,
		Embedder:          embedder,
		RepositoryFetcher: repositoryFetcher,
		TextSplitter:      textSplitter,
		DocumentService:   documentService,
		IndexingConfig:    indexingConfig,
	}
}

// NewIndexRepositoryUseCase creates an index repository use case from the container.
func (c *Container) NewIndexRepositoryUseCase() *indexing.IndexRepositoryUseCase {
	return indexing.NewIndexRepositoryUseCase(
		c.RepositoryFetcher,
		c.VectorStore,
		c.Embedder,
		c.TextSplitter,
		c.IndexingConfig,
	)
}
