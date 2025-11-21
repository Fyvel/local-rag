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

type Service interface {
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)

	CompleteStream(ctx context.Context, req CompletionRequest) (io.ReadCloser, error)

	GenerateTitle(ctx context.Context, question string, model string) (string, error)
}
