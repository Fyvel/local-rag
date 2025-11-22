package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"local-ai/internal/domain/chat"
	"local-ai/internal/infra/httpclient"
)

const (
	DefaultBaseURL = "http://localhost:11434"
	ChatEndpoint   = "/api/chat"
	DefaultModel   = "gpt-oss:20b"
)

// ChatClient implements the chat.Service interface for Ollama.
type ChatClient struct {
	baseURL    string
	httpClient *httpclient.Client
	model      string
}

// ChatOptions holds configuration for the chat client.
type ChatOptions struct {
	BaseURL    string
	HTTPClient *httpclient.Client
	Model      string
}

// ChatOption is a functional option for configuring the chat client.
type ChatOption func(*ChatOptions)

// NewChatClient creates a new Ollama chat client.
func NewChatClient(opts ...ChatOption) chat.Service {
	options := ChatOptions{
		BaseURL:    DefaultBaseURL,
		HTTPClient: httpclient.New(),
		Model:      DefaultModel,
	}

	for _, apply := range opts {
		apply(&options)
	}

	return &ChatClient{
		baseURL:    options.BaseURL,
		httpClient: options.HTTPClient,
		model:      options.Model,
	}
}

// WithChatBaseURL sets the base URL for the Ollama API.
func WithChatBaseURL(baseURL string) ChatOption {
	return func(o *ChatOptions) {
		o.BaseURL = baseURL
	}
}

// WithChatHTTPClient sets the HTTP client.
func WithChatHTTPClient(httpClient *httpclient.Client) ChatOption {
	return func(o *ChatOptions) {
		o.HTTPClient = httpClient
	}
}

// WithChatModel sets the default model to use.
func WithChatModel(model string) ChatOption {
	return func(o *ChatOptions) {
		o.Model = model
	}
}

// ollamaMessage represents a message in Ollama's format.
type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ollamaChatRequest represents Ollama's chat request format.
type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

// ollamaChatResponse represents Ollama's chat response format.
type ollamaChatResponse struct {
	Model     string        `json:"model"`
	CreatedAt string        `json:"created_at"`
	Message   ollamaMessage `json:"message"`
	Done      bool          `json:"done"`
}

// Complete generates a non-streaming completion.
func (c *ChatClient) Complete(ctx context.Context, req chat.CompletionRequest) (*chat.CompletionResponse, error) {
	// Convert messages to Ollama format
	ollamaMessages := make([]ollamaMessage, len(req.Messages))
	for i, msg := range req.Messages {
		ollamaMessages[i] = ollamaMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	model := req.Model
	if model == "" {
		model = c.model
	}

	// Create request
	ollamaReq := ollamaChatRequest{
		Model:    model,
		Messages: ollamaMessages,
		Stream:   false,
	}

	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	url := c.baseURL + ChatEndpoint
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var ollamaResp ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &chat.CompletionResponse{
		Content: ollamaResp.Message.Content,
	}, nil
}

// CompleteStream generates a streaming completion.
func (c *ChatClient) CompleteStream(ctx context.Context, req chat.CompletionRequest) (io.ReadCloser, error) {
	// Convert messages to Ollama format
	ollamaMessages := make([]ollamaMessage, len(req.Messages))
	for i, msg := range req.Messages {
		ollamaMessages[i] = ollamaMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	model := req.Model
	if model == "" {
		model = c.model
	}

	// Create request with streaming enabled
	ollamaReq := ollamaChatRequest{
		Model:    model,
		Messages: ollamaMessages,
		Stream:   true,
	}

	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	url := c.baseURL + ChatEndpoint
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("ollama API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Return a streaming reader that processes the NDJSON stream
	return newStreamReader(resp.Body), nil
}

// streamReader wraps the response body and processes Ollama's NDJSON stream.
type streamReader struct {
	scanner *bufio.Scanner
	body    io.ReadCloser
	buffer  *bytes.Buffer
}

func newStreamReader(body io.ReadCloser) *streamReader {
	return &streamReader{
		scanner: bufio.NewScanner(body),
		body:    body,
		buffer:  &bytes.Buffer{},
	}
}

func (sr *streamReader) Read(p []byte) (n int, err error) {
	// If buffer has data, return it first
	if sr.buffer.Len() > 0 {
		return sr.buffer.Read(p)
	}

	// Read next line from stream
	if !sr.scanner.Scan() {
		if err := sr.scanner.Err(); err != nil {
			return 0, err
		}
		return 0, io.EOF
	}

	line := sr.scanner.Bytes()
	if len(line) == 0 {
		return sr.Read(p) // Skip empty lines
	}

	// Parse the NDJSON line
	var resp ollamaChatResponse
	if err := json.Unmarshal(line, &resp); err != nil {
		return 0, fmt.Errorf("failed to parse stream chunk: %w", err)
	}

	// Write content to buffer
	content := resp.Message.Content
	sr.buffer.WriteString(content)

	// Read from buffer
	return sr.buffer.Read(p)
}

func (sr *streamReader) Close() error {
	return sr.body.Close()
}

func (c *ChatClient) GenerateTitle(ctx context.Context, question string, model string) (string, error) {
	if model == "" {
		model = "mistral:latest"
	}

	systemPrompt := `
		Create a concise, 3 to 5-word phrase (30 characters max) with an emoji as a title for the previous query.
		Do not use the word title.
		Do not use any formatting.

		Examples of titles:
		😢 Sad Story
		🎂 How To Bake A Cake
		✉️ Email Draft
		💻 Programming Help
	`
	messages := []ollamaMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: question,
		},
	}

	ollamaReq := ollamaChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.baseURL + ChatEndpoint
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	title := ollamaResp.Message.Content
	titleBytes := bytes.Trim([]byte(title), `"'`)
	titleBytes = bytes.TrimSpace(titleBytes)
	return string(titleBytes), nil
}
