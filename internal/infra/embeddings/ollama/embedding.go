package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"local-ai/internal/domain/embeddings"
	"local-ai/internal/infra/httpclient"
)

// EmbeddingRequest is serialized and sent to the API server.
type EmbeddingRequest struct {
	Prompt any    `json:"prompt"`
	Model  string `json:"model"`
}

// EmbeddingResponse received from API.
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// ToEmbeddings converts the API response into a slice of embeddings.
func (e *EmbeddingResponse) ToEmbeddings() ([]*embeddings.Embedding, error) {
	return []*embeddings.Embedding{
		{Vector: e.Embedding},
	}, nil
}

// Embed implements the domain embeddings.Embedder interface.
// It generates embeddings for the given text using the specified model.
func (c *Client) Embed(ctx context.Context, text string, model string) ([]*embeddings.Embedding, error) {
	embReq := &EmbeddingRequest{
		Prompt: text,
		Model:  model,
	}

	u, err := url.Parse(c.opts.BaseURL + "/embeddings")
	if err != nil {
		return nil, err
	}

	var body = &bytes.Buffer{}
	enc := json.NewEncoder(body)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(embReq); err != nil {
		return nil, err
	}

	options := []httpclient.RequestOption{}
	req, err := httpclient.NewRequest(ctx, http.MethodPost, u.String(), body, options...)
	if err != nil {
		return nil, err
	}

	resp, err := httpclient.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	e := new(EmbeddingResponse)
	if err := json.NewDecoder(resp.Body).Decode(e); err != nil {
		return nil, err
	}

	return e.ToEmbeddings()
}
