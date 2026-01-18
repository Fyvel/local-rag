# Copilot Instructions — Local RAG API (Go DDD)

This is a Go-based RAG (Retrieval-Augmented Generation) API following **Domain-Driven Design (DDD)** and **Clean Architecture** principles.

## Project Overview

- **Language**: Go 1.23+
- **Framework**: Gin for HTTP routing
- **Architecture**: Hexagonal/Clean Architecture with DDD
- **Key Dependencies**: Ollama (LLM), ChromaDB (vector store), LangChain concepts

## Directory Structure

```
server/           → Entry point (main.go)
internal/
  api/            → HTTP handlers, routing, DTOs
  application/    → Use cases, orchestration logic
  domain/         → Core business logic, entities, interfaces
  infra/          → External integrations, adapters
  di/             → Dependency injection container
```

## Dependency Rules (CRITICAL)

Always enforce strict dependency flow:

### Allowed imports:
- `api` → `application`
- `application` → `domain` (interfaces only)
- `infra` → `domain` (implements interfaces)
- `di` → all layers (wiring only)

### Forbidden imports:
- `domain` → anything (must be pure, zero internal imports)
- `application` → `infra` (must use domain interfaces)
- `api` → `domain` directly (go through application)
- `api` → `infra` (except main/di setup)

## Code Style Guidelines

### Go Idioms
- Use `context.Context` as the first parameter for all I/O operations
- Prefer returning errors over panicking
- Use functional options pattern for configuration (e.g., `WithBaseURL()`)
- Follow "accept interfaces, return structs" principle
- Use small, focused interfaces (Go style)

### Naming Conventions
- Use cases: `*UseCase` suffix (e.g., `AskQuestionUseCase`)
- Commands: `*Command` suffix for input DTOs
- Results: `*Result` suffix for output DTOs
- Handlers: `*Handler` or `*HandlerFactory` suffix
- No package stuttering (use `documents.Document`, not `documents.DocumentDocument`)

### Error Handling
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Use domain-specific errors in the domain layer
- Return structured error responses from API handlers

### Testing
- Write table-driven tests
- Use interfaces for mocking dependencies
- Place test files next to implementation (`*_test.go`)

## DDD Patterns to Follow

### Aggregates
- Each aggregate has a root entity that enforces invariants
- Aggregates are stored/retrieved as a whole
- Reference other aggregates by ID only

### Repositories
- Define repository interfaces in the domain layer
- Implement repositories in the infrastructure layer
- Use domain types in repository methods, not infrastructure types

### Domain Services
- Use when logic doesn't fit in an entity
- Keep services stateless
- Inject through interfaces

### Use Cases (Application Services)
- One use case per business operation
- Accept command objects, return result objects
- Orchestrate domain logic without knowing infrastructure

## Streaming Implementation (SSE/WebSocket)

When implementing streaming endpoints:

### SSE Requirements
```go
// Required headers
w.Header().Set("Content-Type", "text/event-stream")
w.Header().Set("Cache-Control", "no-cache")
w.Header().Set("Connection", "keep-alive")

// Stream format
fmt.Fprintf(w, "data: %s\n\n", chunk)
flusher.Flush()
```

### Streaming Use Cases
- Return `<-chan StreamResponse` for streaming operations
- Use buffered channels for backpressure handling
- Always check `context.Done()` for cancellation
- Close channels properly to signal completion
- Save final messages to repository after stream completes

## API Conventions

### Request/Response
- Use JSON for request/response bodies
- Validate input before processing
- Return appropriate HTTP status codes
- Use consistent error response format

### Endpoint Patterns
```
POST /api/v1/questions           → Quick question (auto-creates discussion)
POST /api/v1/questions/stream    → Streaming quick question
POST /api/v1/discussions         → Create discussion
GET  /api/v1/discussions         → List discussions
GET  /api/v1/discussions/:id     → Get discussion
POST /api/v1/discussions/:id/question        → Ask in discussion
POST /api/v1/discussions/:id/question/stream → Streaming ask
```

## Common Patterns

### Handler Factory Pattern
```go
func HandlerFactory(useCase *UseCase) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Validate input
        // Call use case
        // Return response
    }
}
```

### Functional Options
```go
type Option func(*Config)

func WithBaseURL(url string) Option {
    return func(c *Config) {
        c.BaseURL = url
    }
}
```

### Channel-Based Streaming
```go
func (uc *UseCase) ExecuteStream(ctx context.Context, cmd Command) (<-chan Response, error) {
    out := make(chan Response, 10) // Buffered for backpressure
    go func() {
        defer close(out)
        // Stream processing with ctx.Done() checks
    }()
    return out, nil
}
```

## Build & Run Commands

```bash
make deps        # Install dependencies
make build       # Build the server
make run         # Run the server
make test        # Run tests
swag init -g server/main.go -o docs  # Generate Swagger docs
```
