---
applyTo: "internal/infra/**/*.go"
---

# Infrastructure Layer Instructions

The infrastructure layer provides concrete implementations of domain interfaces. It handles external systems like databases, APIs, and third-party services.

## Critical Rules

### Allowed Imports
- `internal/domain/*` - To implement domain interfaces
- External libraries (Ollama SDK, ChromaDB client, etc.)
- Standard library

### Forbidden Imports
- `internal/api/` - Never import handlers
- `internal/application/` - Never import use cases
- `internal/di/` - Never import container

## Package Organization

```
internal/infra/
  chat/ollama/        → Ollama chat client (implements chat.Service)
  embeddings/ollama/  → Ollama embeddings (implements embeddings.Embedder)
  discussion/memory/  → In-memory repository (implements discussion.Repository)
  store/chroma/       → ChromaDB client
  store/vectorstore/  → Vector store adapter (implements documents.DocumentRepository)
  fetcher/            → GitHub fetcher (implements repositories.Fetcher)
  httpclient/         → HTTP client utilities
  textsplitter/       → Text chunking (implements transformers.TextChunker)
```

## Adapter Pattern

Implement domain interfaces without leaking infrastructure details:
```go
package ollama

import (
    "context"
    "local-ai/internal/domain/chat"
)

// ChatClient implements chat.Service
type ChatClient struct {
    baseURL    string
    httpClient *httpclient.Client
    model      string
}

// Ensure interface compliance at compile time
var _ chat.Service = (*ChatClient)(nil)

func (c *ChatClient) Complete(ctx context.Context, req chat.CompletionRequest) (*chat.CompletionResponse, error) {
    // Transform domain types to Ollama API types
    ollamaReq := toOllamaRequest(req)

    // Make HTTP call
    resp, err := c.httpClient.Post(ctx, c.baseURL+"/api/chat", ollamaReq)
    if err != nil {
        return nil, fmt.Errorf("ollama request failed: %w", err)
    }

    // Transform back to domain types
    return toDomainResponse(resp), nil
}
```

## Functional Options

Use functional options for configuration:
```go
type ChatOption func(*ChatOptions)

type ChatOptions struct {
    BaseURL    string
    HTTPClient *httpclient.Client
    Model      string
}

func NewChatClient(opts ...ChatOption) chat.Service {
    options := ChatOptions{
        BaseURL: DefaultBaseURL,
        Model:   DefaultModel,
    }

    for _, opt := range opts {
        opt(&options)
    }

    return &ChatClient{
        baseURL: options.BaseURL,
        model:   options.Model,
    }
}

func WithChatBaseURL(url string) ChatOption {
    return func(o *ChatOptions) {
        o.BaseURL = url
    }
}

func WithChatModel(model string) ChatOption {
    return func(o *ChatOptions) {
        o.Model = model
    }
}
```

## Streaming Implementation

For streaming responses from external APIs:
```go
func (c *ChatClient) CompleteStreamChannel(ctx context.Context, req chat.CompletionRequest) (<-chan string, error) {
    out := make(chan string, 10) // Buffered for backpressure

    go func() {
        defer close(out)

        resp, err := c.httpClient.Post(ctx, c.baseURL+"/api/chat", req)
        if err != nil {
            return
        }
        defer resp.Body.Close()

        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            select {
            case <-ctx.Done():
                return
            default:
                // Parse and send token
                var chunk ollamaChatResponse
                if err := json.Unmarshal(scanner.Bytes(), &chunk); err == nil {
                    out <- chunk.Message.Content
                }
            }
        }
    }()

    return out, nil
}
```

## Repository Implementation

```go
package memory

import (
    "context"
    "sync"
    "local-ai/internal/domain/discussion"
)

type InMemoryDiscussionRepository struct {
    mu          sync.RWMutex
    discussions map[string]*discussion.Discussion
}

var _ discussion.Repository = (*InMemoryDiscussionRepository)(nil)

func NewInMemoryDiscussionRepository() *InMemoryDiscussionRepository {
    return &InMemoryDiscussionRepository{
        discussions: make(map[string]*discussion.Discussion),
    }
}

func (r *InMemoryDiscussionRepository) FindByID(ctx context.Context, id string) (*discussion.Discussion, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    disc, ok := r.discussions[id]
    if !ok {
        return nil, fmt.Errorf("discussion not found: %s", id)
    }
    return disc, nil
}

func (r *InMemoryDiscussionRepository) Save(ctx context.Context, d *discussion.Discussion) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.discussions[d.ID] = d
    return nil
}
```

## Error Handling

Wrap infrastructure errors appropriately:
```go
func (c *Client) Fetch(ctx context.Context, url string) (*Result, error) {
    resp, err := c.httpClient.Get(ctx, url)
    if err != nil {
        return nil, fmt.Errorf("HTTP request to %s failed: %w", url, err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    return parseResponse(resp)
}
```

## Context Propagation

Always propagate context for cancellation:
```go
func (c *Client) Query(ctx context.Context, query string) ([]Result, error) {
    // Check context early
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }

    // Use context in HTTP calls
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := c.httpClient.Do(req)
    // ...
}
```
