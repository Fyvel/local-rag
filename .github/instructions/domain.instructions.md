---
applyTo: "internal/domain/**/*.go"
---

# Domain Layer Instructions

This is the **pure domain layer** - the heart of the application. It must have ZERO dependencies on other internal packages.

## Critical Rules

### No External Imports
The domain layer must NOT import from:
- `internal/api/`
- `internal/application/`
- `internal/infra/`
- `internal/di/`

Only standard library and domain packages within `internal/domain/` are allowed.

### Package Structure
Each domain package represents a bounded context or aggregate:
- `discussion/` - Discussion aggregate (messages, sessions)
- `documents/` - Document aggregate (content, metadata)
- `embeddings/` - Embedding value objects and interfaces
- `repositories/` - Repository aggregate (GitHub repos)
- `transformers/` - Text processing interfaces
- `chat/` - Chat service interface

## Entity Design

### Aggregate Roots
```go
// Aggregate root with invariant enforcement
type Discussion struct {
    ID        string
    Title     string
    Messages  []*Message
    Status    DiscussionStatus
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Factory function that validates invariants
func NewDiscussion(id, title string) (*Discussion, error) {
    if id == "" {
        return nil, errors.New("discussion ID cannot be empty")
    }
    // ... enforce all invariants
}
```

### Value Objects
- Immutable after creation
- Equality based on value, not identity
- No setters, only constructors

### Domain Errors
Define domain-specific errors:
```go
var (
    ErrDiscussionClosed = errors.New("discussion is closed")
    ErrEmptyMessage     = errors.New("message cannot be empty")
)
```

## Interface Definitions

Define interfaces that infrastructure will implement:
```go
// Repository interface - implemented by infra
type Repository interface {
    FindByID(ctx context.Context, id string) (*Discussion, error)
    Save(ctx context.Context, d *Discussion) error
    List(ctx context.Context) ([]*Discussion, error)
}

// Service interface - implemented by infra
type Embedder interface {
    Embed(ctx context.Context, text string) ([]float64, error)
}
```

## Domain Services

Use domain services for logic that doesn't fit in entities:
```go
type StreamingService interface {
    StreamTokens(ctx context.Context, discussionID string) (<-chan string, error)
}
```

## Invariants

Always enforce business rules in the domain:
```go
func (d *Discussion) AddMessage(msg *Message) error {
    if d.Status == StatusClosed {
        return ErrDiscussionClosed
    }
    if msg.Content == "" {
        return ErrEmptyMessage
    }
    d.Messages = append(d.Messages, msg)
    d.UpdatedAt = time.Now()
    return nil
}
```
