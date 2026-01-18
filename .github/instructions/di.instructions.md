---
applyTo: "internal/di/**/*.go"
---

# Dependency Injection Layer Instructions

The DI layer wires all application dependencies together. It's the only layer allowed to import from all other layers.

## Critical Rules

### Allowed Imports
- `internal/domain/*` - For interface types
- `internal/application/*` - For use case constructors
- `internal/infra/*` - For concrete implementations
- Standard library

### Purpose
- Create and configure infrastructure components
- Wire domain interfaces to infrastructure implementations
- Construct use cases with their dependencies
- Provide factories for the API layer

## Container Structure

```go
package di

import (
    "local-ai/internal/domain/chat"
    "local-ai/internal/domain/discussion"
    "local-ai/internal/domain/documents"
    "local-ai/internal/domain/embeddings"
    
    appdiscussion "local-ai/internal/application/discussion"
    "local-ai/internal/application/indexing"
    
    "local-ai/internal/infra/chat/ollama"
    "local-ai/internal/infra/discussion/memory"
    infraollama "local-ai/internal/infra/embeddings/ollama"
)

// Config holds external configuration
type Config struct {
    OllamaURL  string
    ChromaURL  string
    ChatModel  string
    TitleModel string
}

// Container holds all wired dependencies
type Container struct {
    // Infrastructure components
    HTTPClient           *httpclient.Client
    VectorStore          documents.DocumentRepository
    Embedder             embeddings.Embedder
    RepositoryFetcher    repositories.Fetcher
    TextSplitter         transformers.TextChunker
    DiscussionRepository discussion.Repository
    ChatService          chat.Service

    // Configuration
    IndexingConfig indexing.Config
    TitleModel     string
}

// NewContainer creates and wires all dependencies
func NewContainer(cfg Config) *Container {
    // 1. Create shared infrastructure
    httpClient := httpclient.New()

    // 2. Create domain interface implementations
    embedder := infraollama.NewClient(
        infraollama.WithBaseURL(cfg.OllamaURL),
        infraollama.WithHTTPClient(httpClient),
    )

    chatService := ollama.NewChatClient(
        ollama.WithChatBaseURL(cfg.OllamaURL),
        ollama.WithChatHTTPClient(httpClient),
        ollama.WithChatModel(cfg.ChatModel),
    )

    discussionRepo := memory.NewInMemoryDiscussionRepository()

    // 3. Return wired container
    return &Container{
        HTTPClient:           httpClient,
        Embedder:             embedder,
        ChatService:          chatService,
        DiscussionRepository: discussionRepo,
        // ... other dependencies
    }
}
```

## Use Case Factories

Create factory methods for each use case:
```go
// Each factory method constructs a use case with its dependencies
func (c *Container) NewAskQuestionUseCase() *appdiscussion.AskQuestionUseCase {
    return appdiscussion.NewAskQuestionUseCase(
        c.DiscussionRepository,
        c.ChatService,
    )
}

func (c *Container) NewAskQuestionStreamUseCase() *appdiscussion.AskQuestionStreamUseCase {
    return appdiscussion.NewAskQuestionStreamUseCase(
        c.DiscussionRepository,
        c.ChatService,
    )
}

func (c *Container) NewCreateDiscussionUseCase() *appdiscussion.CreateDiscussionUseCase {
    return appdiscussion.NewCreateDiscussionUseCase(
        c.DiscussionRepository,
    )
}

func (c *Container) NewListDiscussionsUseCase() *appdiscussion.ListDiscussionsUseCase {
    return appdiscussion.NewListDiscussionsUseCase(
        c.DiscussionRepository,
    )
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
```

## Configuration Defaults

Handle missing configuration gracefully:
```go
func NewContainer(cfg Config) *Container {
    // Set defaults for optional config
    chatModel := cfg.ChatModel
    if chatModel == "" {
        chatModel = ollama.DefaultModel
    }

    ollamaURL := cfg.OllamaURL
    if ollamaURL == "" {
        ollamaURL = ollama.DefaultBaseURL
    }

    // ...
}
```

## Best Practices

1. **Single Responsibility**: Container only wires, no business logic
2. **Explicit Dependencies**: All dependencies visible in constructor calls
3. **Interface Types**: Store domain interfaces, not concrete types
4. **Lazy Initialization**: Create use cases on demand via factory methods
5. **Configuration Validation**: Validate required config at container creation
