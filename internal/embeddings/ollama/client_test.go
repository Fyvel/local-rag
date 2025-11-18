package ollama

import (
	"local-ai/internal/client"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Run("BaseURL", func(t *testing.T) {
		c := NewClient()
		assert.Equal(t, c.opts.BaseURL, BaseURL)

		testVal := "http://foo"
		c = NewClient(WithBaseURL(testVal))
		assert.Equal(t, c.opts.BaseURL, testVal)
	})

	t.Run("http client", func(t *testing.T) {
		c := NewClient()
		assert.NotNil(t, c.opts.HTTPClient)

		testVal := client.NewHTTP()
		c = NewClient(WithHTTPClient(testVal))
		assert.NotNil(t, c.opts.HTTPClient)
	})
}
