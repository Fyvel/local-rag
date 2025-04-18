package chroma

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type AddDocumentsRequest struct {
	IDs        []string         `json:"ids"`
	Embeddings [][]float64      `json:"embeddings"`
	Metadatas  []map[string]any `json:"metadatas,omitempty"`
	Documents  []string         `json:"documents,omitempty"`
}

func (c *Client) AddDocuments(ctx context.Context, collectionID string, req *AddDocumentsRequest) error {
	url := fmt.Sprintf("%s/api/v2/tenants/%s/databases/%s/collections/%s/add",
		c.baseURL, c.tenant, c.database, collectionID)

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		httpReq.Header.Set("x-chroma-token", c.token)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to add documents: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add documents, status: %d", resp.StatusCode)
	}

	return nil
}

type DeleteDocumentsRequest struct {
	IDs []string `json:"ids"`
}

func (c *Client) DeleteDocuments(ctx context.Context, collectionID string, ids []string) error {
	url := fmt.Sprintf("%s/api/v2/tenants/%s/databases/%s/collections/%s/delete",
		c.baseURL, c.tenant, c.database, collectionID)

	req := DeleteDocumentsRequest{
		IDs: ids,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		httpReq.Header.Set("x-chroma-token", c.token)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to delete documents: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete documents, status: %d", resp.StatusCode)
	}

	return nil
}

type QueryRequest struct {
	QueryEmbeddings [][]float64 `json:"query_embeddings"`
	NResults        int         `json:"n_results,omitempty"`
	Include         []string    `json:"include,omitempty"`
}

type QueryResponse struct {
	IDs        [][]string         `json:"ids"`
	Distances  [][]float64        `json:"distances,omitempty"`
	Documents  [][]string         `json:"documents,omitempty"`
	Metadatas  [][]map[string]any `json:"metadatas,omitempty"`
	Embeddings [][][]float64      `json:"embeddings,omitempty"`
	Include    []string           `json:"include"`
}

type QueryResult struct {
	ID       string
	Distance float64
	Document string
	Metadata map[string]any
}

func (c *Client) Query(ctx context.Context, collectionID string, embedding []float64, limit int) ([]QueryResult, error) {
	url := fmt.Sprintf("%s/api/v2/tenants/%s/databases/%s/collections/%s/query",
		c.baseURL, c.tenant, c.database, collectionID)

	req := QueryRequest{
		QueryEmbeddings: [][]float64{embedding},
		NResults:        limit,
		Include:         []string{"documents", "metadatas", "distances"},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		httpReq.Header.Set("x-chroma-token", c.token)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to query documents: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to query documents, status: %d", resp.StatusCode)
	}

	var queryResp QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&queryResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var results []QueryResult

	if len(queryResp.IDs) > 0 {
		for i, idList := range queryResp.IDs {
			for j, id := range idList {
				result := QueryResult{
					ID: id,
				}

				if len(queryResp.Distances) > i && len(queryResp.Distances[i]) > j {
					result.Distance = queryResp.Distances[i][j]
				}

				if len(queryResp.Documents) > i && len(queryResp.Documents[i]) > j {
					result.Document = queryResp.Documents[i][j]
				}

				if len(queryResp.Metadatas) > i && len(queryResp.Metadatas[i]) > j {
					result.Metadata = queryResp.Metadatas[i][j]
				}

				results = append(results, result)
			}
		}
	}

	return results, nil
}
