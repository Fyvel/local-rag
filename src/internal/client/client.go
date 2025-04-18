package client

import (
	"context"
	"net/http"
)

type HTTP struct {
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

func NewHTTP(opts ...Option) *HTTP {
	options := Options{
		HTTPClient: &http.Client{},
	}
	for _, apply := range opts {
		apply(&options)
	}

	return &HTTP{
		client:  options.HTTPClient,
		limiter: options.Limiter,
	}
}

func (h *HTTP) Do(req *http.Request) (*http.Response, error) {
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
