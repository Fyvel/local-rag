package ollama

import (
	"local-ai/internal/domain/embeddings"
	"local-ai/internal/infra/httpclient"
)

const (
	BaseURL = "http://localhost:11434/api"
)

type Client struct {
	opts Options
}

type Options struct {
	BaseURL    string
	HTTPClient *httpclient.Client
}

type Option func(*Options)

func NewClient(opts ...Option) *Client {
	options := Options{
		BaseURL:    BaseURL,
		HTTPClient: httpclient.New(),
	}

	for _, apply := range opts {
		apply(&options)
	}

	return &Client{
		opts: options,
	}
}

// NewEmbedder creates a new Ollama embedder that implements the domain Embedder interface.
func NewEmbedder(opts ...Option) embeddings.Embedder {
	return NewClient(opts...)
}

func WithBaseURL(baseURL string) Option {
	return func(o *Options) {
		o.BaseURL = baseURL
	}
}

func WithHTTPClient(httpClient *httpclient.Client) Option {
	return func(o *Options) {
		o.HTTPClient = httpClient
	}
}
