package ollama

import (
	"indexing/client"
	"indexing/internal/embeddings"
)

const (
	BaseURL = "http://localhost:11434/api"
)

type Client struct {
	opts Options
}

type Options struct {
	BaseURL    string
	HTTPClient *client.HTTP
}

type Option func(*Options)

func NewClient(opts ...Option) *Client {
	options := Options{
		BaseURL:    BaseURL,
		HTTPClient: client.NewHTTP(),
	}

	for _, apply := range opts {
		apply(&options)
	}

	return &Client{
		opts: options,
	}
}

func NewEmbedder(opts ...Option) embeddings.Embedder[*EmbeddingRequest] {
	return NewClient(opts...)
}

func WithBaseURL(baseURL string) Option {
	return func(o *Options) {
		o.BaseURL = baseURL
	}
}

func WithHTTPClient(httpClient *client.HTTP) Option {
	return func(o *Options) {
		o.HTTPClient = httpClient
	}
}
