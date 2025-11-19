package httpclient

import (
	"context"
	"net/http"
)

// Client wraps an HTTP client with optional rate limiting.
type Client struct {
	client  *http.Client
	limiter Limiter
}

type Options struct {
	HTTPClient *http.Client
	Limiter    Limiter
}

type Option func(*Options)

type Limiter interface {
	Wait(context.Context) error
}

func New(opts ...Option) *Client {
	options := Options{
		HTTPClient: &http.Client{},
	}
	for _, apply := range opts {
		apply(&options)
	}

	return &Client{
		client:  options.HTTPClient,
		limiter: options.Limiter,
	}
}

func (h *Client) Do(req *http.Request) (*http.Response, error) {
	if h.limiter != nil {
		err := h.limiter.Wait(req.Context())
		if err != nil {
			return nil, err
		}
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func WithHTTPClient(c *http.Client) Option {
	return func(o *Options) {
		o.HTTPClient = c
	}
}

func WithLimiter(l Limiter) Option {
	return func(o *Options) {
		o.Limiter = l
	}
}
