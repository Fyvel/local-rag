package discussion

import (
	"context"
)

type StreamToken struct {
	Content string
	Done    bool
	Error   error
}

type StreamingService interface {
	StreamAnswer(ctx context.Context, disc *Discussion, question string) (<-chan StreamToken, error)
}
