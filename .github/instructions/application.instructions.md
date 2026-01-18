---
applyTo: "internal/application/**/*.go"
---

# Application Layer Instructions

The application layer contains **use cases** that orchestrate domain logic. It depends only on domain interfaces.

## Critical Rules

### Allowed Imports
- `internal/domain/*` (interfaces only)
- Standard library

### Forbidden Imports
- `internal/api/` - Never import handlers
- `internal/infra/` - Never import concrete implementations
- `internal/di/` - Never import container

## Use Case Structure

### Naming Convention
- One file per use case: `ask_question.go`, `create_discussion.go`
- Struct suffix: `*UseCase`
- Constructor: `New*UseCase`

### Standard Pattern
```go
package discussion

import (
    "context"
    "local-ai/internal/domain/discussion"
    "local-ai/internal/domain/chat"
)

// Command - input DTO
type AskQuestionCommand struct {
    DiscussionID string
    Question     string
}

// Result - output DTO
type AskQuestionResult struct {
    MessageID string
    Answer    string
}

// UseCase struct with domain interface dependencies
type AskQuestionUseCase struct {
    repo        discussion.Repository  // Domain interface
    chatService chat.Service           // Domain interface
}

// Constructor with dependency injection
func NewAskQuestionUseCase(
    repo discussion.Repository,
    chatService chat.Service,
) *AskQuestionUseCase {
    return &AskQuestionUseCase{
        repo:        repo,
        chatService: chatService,
    }
}

// Execute method - the use case logic
func (uc *AskQuestionUseCase) Execute(
    ctx context.Context,
    cmd AskQuestionCommand,
) (*AskQuestionResult, error) {
    // 1. Validate command
    if cmd.DiscussionID == "" {
        return nil, fmt.Errorf("discussion ID required")
    }

    // 2. Load aggregate
    disc, err := uc.repo.FindByID(ctx, cmd.DiscussionID)
    if err != nil {
        return nil, fmt.Errorf("failed to load discussion: %w", err)
    }

    // 3. Execute domain logic
    // ...

    // 4. Persist changes
    if err := uc.repo.Save(ctx, disc); err != nil {
        return nil, fmt.Errorf("failed to save: %w", err)
    }

    // 5. Return result
    return &AskQuestionResult{...}, nil
}
```

## Streaming Use Cases

For streaming operations, return channels:
```go
type StreamResponse struct {
    MessageID string `json:"message_id,omitempty"`
    Token     string `json:"token,omitempty"`
    Done      bool   `json:"done"`
    Error     string `json:"error,omitempty"`
}

func (uc *AskQuestionStreamUseCase) Execute(
    ctx context.Context,
    cmd AskQuestionCommand,
) (<-chan StreamResponse, error) {
    // Validate before starting stream
    if cmd.Question == "" {
        return nil, fmt.Errorf("question cannot be empty")
    }

    // Create buffered channel
    out := make(chan StreamResponse, 10)

    go func() {
        defer close(out)

        // Stream tokens
        for token := range tokenChan {
            select {
            case <-ctx.Done():
                return
            case out <- StreamResponse{Token: token}:
            }
        }

        // Send completion
        out <- StreamResponse{Done: true}
    }()

    return out, nil
}
```

## Context Handling

Always respect context cancellation:
```go
func (uc *UseCase) Execute(ctx context.Context, cmd Command) error {
    // Check context before expensive operations
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // Use context with timeouts
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    // Pass context to all operations
    result, err := uc.repo.FindByID(ctx, cmd.ID)
}
```

## Error Handling

Wrap errors with context:
```go
if err != nil {
    return nil, fmt.Errorf("failed to retrieve discussion %s: %w", cmd.ID, err)
}
```
