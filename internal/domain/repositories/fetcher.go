package repositories

import "context"

// Fetcher defines the contract for fetching repositories from external sources.
// This is a domain interface that abstracts away the specific implementation details
// of how repositories are retrieved (e.g., from GitHub, GitLab, local filesystem).
type Fetcher interface {
	// FetchRepository retrieves a repository and its files based on the given URL and filters.
	// Returns a Repository aggregate containing all files matching the specified file types.
	FetchRepository(ctx context.Context, url string, fileTypes []string) (*Repository, error)
}
