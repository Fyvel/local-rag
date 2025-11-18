package chroma

import (
	"local-ai/internal/client"
)

type Client struct {
	httpClient *client.HTTP
	baseURL    string
	tenant     string
	database   string
	collection string
	token      string
}

type Option func(*Client)

func WithChromaURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

func WithHTTPClient(httpClient *client.HTTP) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithCollection(collection string) Option {
	return func(c *Client) {
		c.collection = collection
	}
}

func WithTenant(tenant string) Option {
	return func(c *Client) {
		c.tenant = tenant
	}
}

func WithDatabase(database string) Option {
	return func(c *Client) {
		c.database = database
	}
}

func WithToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL:  "http://localhost:8000",
		tenant:   "default_tenant",
		database: "default_database",
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = client.NewHTTP()
	}

	return c
}
