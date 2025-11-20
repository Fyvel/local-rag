package chroma

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Collection struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
	Tenant   string                 `json:"tenant"`
	Database string                 `json:"database"`
}

func (c *Client) GetCollections(ctx context.Context) ([]Collection, error) {
	url := fmt.Sprintf("%s/api/v2/tenants/%s/databases/%s/collections",
		c.baseURL, c.tenant, c.database)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("x-chroma-token", c.token)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list collections, status: %d", res.StatusCode)
	}

	var collections []Collection
	if err := json.NewDecoder(res.Body).Decode(&collections); err != nil {
		return nil, fmt.Errorf("failed to decode collections: %w", err)
	}

	return collections, nil
}

func (c *Client) GetCollection(ctx context.Context, name string) (*Collection, error) {
	url := fmt.Sprintf("%s/api/v2/tenants/%s/databases/%s/collections/%s",
		c.baseURL, c.tenant, c.database, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("x-chroma-token", c.token)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("collection not found, status: %d", res.StatusCode)
	}
	var collection Collection
	if err := json.NewDecoder(res.Body).Decode(&collection); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &collection, nil
}

func (c *Client) CreateCollection(ctx context.Context, name string) (*Collection, error) {
	url := fmt.Sprintf("%s/api/v2/tenants/%s/databases/%s/collections",
		c.baseURL, c.tenant, c.database)

	payload := map[string]interface{}{
		"name":          name,
		"get_or_create": true,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("x-chroma-token", c.token)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create collection, status: %d", res.StatusCode)
	}

	var collection Collection
	if err := json.NewDecoder(res.Body).Decode(&collection); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &collection, nil
}

func (c *Client) EnsureCollection(ctx context.Context) (*Collection, error) {
	collection, err := c.GetCollection(ctx, c.collection)
	if err == nil {
		return collection, nil
	}

	return c.CreateCollection(ctx, c.collection)
}
