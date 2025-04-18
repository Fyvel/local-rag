package main

import (
	"context"
	"fmt"
	"indexing/internal/indexer"
)

func main() {
	ctx := context.Background()
	repo, err := indexer.IndexGithubRepository(ctx, "https://github.com/golang/go", []string{"md", "mdx"}, "test_index")
	// repo, err := indexer.IndexGithubRepository(ctx, "https://github.com/Fyvel/local-rag", []string{"md", "mdx"}, "test_index")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, file := range repo.Files {
		if i < 5 { // Print just first 5 files
			fmt.Printf("- %s (%d bytes)\n", file.Path, len(file.Content))
		}
	}
}
