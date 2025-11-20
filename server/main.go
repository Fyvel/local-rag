package main

import (
	"fmt"

	"local-ai/internal/api"
	"local-ai/internal/di"

	_ "local-ai/docs" // This will import the auto-generated docs
)

func main() {
	// Initialize dependency injection container with configuration
	cfg := di.Config{
		ChromaURL:   "http://localhost:8000",
		OllamaURL:   "http://localhost:11434/api",
		TargetIndex: "default_index",
	}
	container := di.NewContainer(cfg)

	// Create router with injected dependencies
	r := api.NewRouter(api.RouterConfig{
		Container: container,
	})

	fmt.Println("Server starting on :8080")
	r.Run(":8080")
}
