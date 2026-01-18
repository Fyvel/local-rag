package chat

import (
	"context"
	"io"
)

type Message struct {
	Role    string
	Content string
}

type CompletionRequest struct {
	Model    string
	Messages []Message
	Stream   bool
}

type CompletionResponse struct {
	Content string
}

type StreamChunk struct {
	Content string
	Done    bool
	Error   error
}

type Service interface {
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)

	CompleteStream(ctx context.Context, req CompletionRequest) (io.ReadCloser, error)

	CompleteStreamChannel(ctx context.Context, req CompletionRequest) (<-chan StreamChunk, error)

	GenerateTitle(ctx context.Context, question string, model string) (string, error)
}
