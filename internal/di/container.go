package di

import (
	"local-ai/internal/application/discussion"
	"local-ai/internal/application/indexing"
	"local-ai/internal/domain/chat"
	domaindiscussion "local-ai/internal/domain/discussion"
	"local-ai/internal/domain/documents"
	"local-ai/internal/domain/embeddings"
	"local-ai/internal/domain/repositories"
	"local-ai/internal/domain/transformers"
	ollamachat "local-ai/internal/infra/chat/ollama"
	"local-ai/internal/infra/discussion/memory"
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
	ChatModel    string
	TitleModel   string // Model used for generating discussion titles
}

// Container holds all application dependencies.
type Container struct {
	HTTPClient           *httpclient.Client
	VectorStore          documents.DocumentRepository
	Embedder             embeddings.Embedder
	RepositoryFetcher    repositories.Fetcher
	TextSplitter         transformers.TextChunker
	IndexingConfig       indexing.Config
	DiscussionRepository domaindiscussion.Repository
	ChatService          chat.Service
	TitleModel           string
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
	embedder := ollama.NewClient(
		ollama.WithBaseURL(cfg.OllamaURL),
		ollama.WithHTTPClient(httpClient),
	)

	// Initialize repository fetcher
	repositoryFetcher := fetcher.NewGitHubFetcher()

	// Initialize indexing config
	indexingConfig := indexing.DefaultConfig()

	// Initialize discussion repository (in-memory for now)
	discussionRepo := memory.NewInMemoryDiscussionRepository()

	// Initialize chat service
	chatModel := cfg.ChatModel
	if chatModel == "" {
		chatModel = "gpt-oss:20b"
	}
	chatService := ollamachat.NewChatClient(
		ollamachat.WithChatBaseURL(cfg.OllamaURL),
		ollamachat.WithChatHTTPClient(httpClient),
		ollamachat.WithChatModel(chatModel),
	)

	// Set title model
	titleModel := cfg.TitleModel
	if titleModel == "" {
		titleModel = "mistral:latest"
	}

	return &Container{
		HTTPClient:           httpClient,
		VectorStore:          vectorStore,
		Embedder:             embedder,
		RepositoryFetcher:    repositoryFetcher,
		TextSplitter:         textSplitter,
		IndexingConfig:       indexingConfig,
		DiscussionRepository: discussionRepo,
		ChatService:          chatService,
		TitleModel:           titleModel,
	}
}

func (c *Container) NewIndexRepositoryUseCase() *indexing.IndexRepositoryUseCase {
	return indexing.NewIndexRepositoryUseCase(
		c.RepositoryFetcher,
		c.VectorStore,
		c.Embedder,
		c.TextSplitter,
		c.IndexingConfig,
	)
}

func (c *Container) NewListDiscussionsUseCase() *discussion.ListDiscussionsUseCase {
	return discussion.NewListDiscussionsUseCase(c.DiscussionRepository)
}

func (c *Container) NewCreateDiscussionUseCase() *discussion.CreateDiscussionUseCase {
	return discussion.NewCreateDiscussionUseCase(c.DiscussionRepository)
}

func (c *Container) NewGetDiscussionUseCase() *discussion.GetDiscussionUseCase {
	return discussion.NewGetDiscussionUseCase(c.DiscussionRepository)
}

func (c *Container) NewAskQuestionUseCase() *discussion.AskQuestionUseCase {
	return discussion.NewAskQuestionUseCase(c.DiscussionRepository, c.ChatService)
}

func (c *Container) NewGetDiscussionHistoryUseCase() *discussion.GetDiscussionHistoryUseCase {
	return discussion.NewGetDiscussionHistoryUseCase(c.DiscussionRepository)
}

func (c *Container) NewAskQuestionQuickUseCase() *discussion.AskQuestionQuickUseCase {
	return discussion.NewAskQuestionQuickUseCase(c.DiscussionRepository, c.ChatService, c.TitleModel)
}
