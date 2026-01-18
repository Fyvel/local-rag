# Agent Instructions — Local RAG API

## Quick Reference

This is a **Go DDD (Domain-Driven Design)** project for a RAG API with streaming support.

### Key Commands
```bash
make build       # Build the application
make run         # Run the server
make test        # Run tests
go build ./...   # Verify compilation
```

### Architecture at a Glance

```
server/main.go        → Entry point
internal/
  di/container.go     → Dependency injection (wires everything)
  api/router.go       → HTTP routing with Gin
  api/*.go            → HTTP handlers
  application/        → Use cases (business orchestration)
  domain/             → Pure business logic (NO external imports!)
  infra/              → External integrations (Ollama, ChromaDB)
```

### Golden Rules

1. **Domain layer is PURE** — No imports from api, application, infra, or di
2. **Application uses domain interfaces** — Never concrete infra types
3. **API calls application** — Never domain or infra directly
4. **Use `context.Context`** — First parameter for all I/O operations
5. **Wrap errors** — `fmt.Errorf("context: %w", err)`

### Common Tasks

#### Adding a New Use Case
1. Create `internal/application/<domain>/<action>.go`
2. Define Command and Result structs
3. Create UseCase struct with domain interface dependencies
4. Add factory method to `internal/di/container.go`
5. Create handler in `internal/api/`
6. Register route in `internal/api/router.go`

#### Adding a Streaming Endpoint
1. Use case returns `<-chan StreamResponse`
2. Handler sets SSE headers
3. Loop over channel, flush after each write
4. Check `ctx.Done()` for cancellation

#### Adding a Domain Entity
1. Create in `internal/domain/<aggregate>/`
2. Use factory function: `NewEntity()`
3. Enforce invariants in methods
4. Define interface if needed by other layers

### Streaming Response Format
```go
type StreamResponse struct {
    MessageID string `json:"message_id,omitempty"`
    Token     string `json:"token,omitempty"`
    Done      bool   `json:"done"`
    Error     string `json:"error,omitempty"`
}
```

### Test Pattern
```go
func TestUseCase_Execute(t *testing.T) {
    tests := []struct {
        name    string
        cmd     Command
        wantErr bool
    }{
        // test cases
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // arrange, act, assert
        })
    }
}
```
