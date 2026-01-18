---
applyTo: "internal/api/**/*.go"
---

# API Layer Instructions

The API layer handles HTTP concerns: routing, request/response handling, and validation. It delegates business logic to the application layer.

## Critical Rules

### Allowed Imports
- `internal/application/*` - Use cases only
- Gin framework (`github.com/gin-gonic/gin`)
- Standard library

### Forbidden Imports
- `internal/domain/` - Never access domain directly (go through application)
- `internal/infra/` - Never access infrastructure directly

## Handler Pattern

### Handler Factories
Use factory functions to inject dependencies:
```go
func AskQuestionHandlerFactory(useCase *discussion.AskQuestionUseCase) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Parse and validate request
        var req AskQuestionRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
            return
        }

        // 2. Extract path parameters
        discussionID := c.Param("id")
        if discussionID == "" {
            c.JSON(http.StatusBadRequest, ErrorResponse{Error: "discussion ID required"})
            return
        }

        // 3. Build command
        cmd := discussion.AskQuestionCommand{
            DiscussionID: discussionID,
            Question:     req.Question,
        }

        // 4. Execute use case
        result, err := useCase.Execute(c.Request.Context(), cmd)
        if err != nil {
            c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
            return
        }

        // 5. Return response
        c.JSON(http.StatusOK, mapToResponse(result))
    }
}
```

## Request/Response Types

Define DTOs in `types.go`:
```go
// Request DTOs
type AskQuestionRequest struct {
    Question string `json:"question" binding:"required"`
}

type CreateDiscussionRequest struct {
    Title string `json:"title" binding:"required"`
}

// Response DTOs
type DiscussionResponse struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}
```

## SSE Streaming Handlers

For Server-Sent Events streaming:
```go
func StreamHandlerFactory(useCase *discussion.AskQuestionStreamUseCase) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Validate request first
        var req AskQuestionRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
            return
        }

        // Set SSE headers
        c.Header("Content-Type", "text/event-stream")
        c.Header("Cache-Control", "no-cache")
        c.Header("Connection", "keep-alive")

        // Get flusher
        flusher, ok := c.Writer.(http.Flusher)
        if !ok {
            c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "streaming not supported"})
            return
        }

        // Execute streaming use case
        stream, err := useCase.Execute(c.Request.Context(), cmd)
        if err != nil {
            c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
            return
        }

        // Stream responses
        for resp := range stream {
            data, _ := json.Marshal(resp)
            fmt.Fprintf(c.Writer, "data: %s\n\n", data)
            flusher.Flush()
        }
    }
}
```

## Validation

Use Gin's binding for validation:
```go
type CreateRequest struct {
    Title    string `json:"title" binding:"required,min=1,max=255"`
    Content  string `json:"content" binding:"required"`
    Priority int    `json:"priority" binding:"gte=1,lte=5"`
}
```

## Router Configuration

In `router.go`:
```go
func NewRouter(cfg RouterConfig) *gin.Engine {
    r := gin.Default()

    // Create use cases from DI container
    askUseCase := cfg.Container.NewAskQuestionUseCase()

    v1 := r.Group("/api/v1")
    {
        // Health check
        v1.GET("/health", HealthCheckHandler)

        // Questions
        v1.POST("/questions", AskQuestionQuickHandlerFactory(quickUseCase))
        v1.POST("/questions/stream", AskQuestionQuickStreamHandlerFactory(streamUseCase))

        // Discussions
        v1.GET("/discussions", ListDiscussionsHandlerFactory(listUseCase))
        v1.POST("/discussions", CreateDiscussionHandlerFactory(createUseCase))
        v1.GET("/discussions/:id", GetDiscussionHandlerFactory(getUseCase))
    }

    return r
}
```

## Error Handling

Use consistent error responses:
```go
func handleError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, discussion.ErrNotFound):
        c.JSON(http.StatusNotFound, ErrorResponse{Error: "not found"})
    case errors.Is(err, discussion.ErrValidation):
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
    default:
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
    }
}
```
