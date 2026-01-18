package api

// =============================================================================
// Standard Response Wrapper
// =============================================================================

// Response is the standard API response format for all endpoints.
// Uses a flat structure with status and result/error fields.
type Response[T any] struct {
	Status string `json:"status" example:"success"` // "success" or "error"
	Result T      `json:"result,omitempty"`         // The response data (omitted on error)
	Error  string `json:"error,omitempty"`          // Error message (omitted on success)
}

// NewSuccessResponse creates a successful response with the given result.
func NewSuccessResponse[T any](result T) Response[T] {
	return Response[T]{
		Status: "success",
		Result: result,
	}
}

// NewErrorResponse creates an error response with the given message.
func NewErrorResponse(err string) Response[any] {
	return Response[any]{
		Status: "error",
		Error:  err,
	}
}

// =============================================================================
// Discussion Request DTOs
// =============================================================================

// CreateDiscussionRequest is the request body for creating a new discussion.
type CreateDiscussionRequest struct {
	Title string `json:"title" binding:"required" example:"How to use embeddings?"`
}

// AskQuestionRequest is the request body for asking a question in a discussion.
type AskQuestionRequest struct {
	Question string `json:"question" binding:"required" example:"What are vector embeddings?"`
}

// =============================================================================
// Indexing Request DTOs
// =============================================================================

// IndexGithubRepoRequest is the request body for indexing a GitHub repository.
type IndexGithubRepoRequest struct {
	GithubURL   string   `json:"github_url" binding:"required" example:"https://github.com/username/repo"`
	FileTypes   []string `json:"file_types" example:"md,mdx"`
	TargetIndex string   `json:"target_index" example:"target_index"`
}
