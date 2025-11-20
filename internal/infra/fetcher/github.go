package fetcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"local-ai/internal/domain/repositories"
)

// GitHubFetcher is an infrastructure implementation of repositories.Fetcher
// that retrieves repositories from GitHub using git operations.
type GitHubFetcher struct{}

// NewGitHubFetcher creates a new GitHub fetcher.
func NewGitHubFetcher() repositories.Fetcher {
	return &GitHubFetcher{}
}

// FetchRepository implements the repositories.Fetcher interface for GitHub.
func (f *GitHubFetcher) FetchRepository(ctx context.Context, githubURL string, fileTypes []string) (*repositories.Repository, error) {
	var repoDir string
	var repo *git.Repository
	var err error

	debugMode := os.Getenv("DEBUG") == "true"
	if debugMode {
		// use persistent directory structure
		segments := strings.Split(strings.TrimSuffix(githubURL, ".git"), "/")
		repoDir = filepath.Join("/tmp/repos", segments[3], segments[4])

		if err = os.MkdirAll(repoDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}
		fmt.Println("DEBUG: Using persistent repository directory:", repoDir)
	} else {
		// use temporary directory that gets cleaned up
		if repoDir, err = os.MkdirTemp("", "github-repo-*"); err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer os.RemoveAll(repoDir)
	}

	// Check if the repository already exists
	if _, err := os.Stat(filepath.Join(repoDir, ".git")); err == nil {
		fmt.Println("Using existing repository at", repoDir)
		repo, err = git.PlainOpen(repoDir)
		if err != nil {
			return nil, fmt.Errorf("failed to open existing repository: %w", err)
		}
	} else {
		fmt.Println("Clone repository started", githubURL)
		repo, err = git.PlainClone(repoDir, false, &git.CloneOptions{
			URL:      githubURL,
			Progress: os.Stdout,
		})
		fmt.Println("Clone operation completed")
		if err != nil {
			return nil, fmt.Errorf("failed to clone repository: %w", err)
		}
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	commitRef, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object: %w", err)
	}

	tree, err := commitRef.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit tree: %w", err)
	}

	// Walk the tree and collect files with matching file types
	var files []repositories.File
	err = tree.Files().ForEach(func(f *object.File) error {
		if shouldIncludeFile(f.Name, fileTypes) {
			content, err := f.Contents()
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", f.Name, err)
			}

			files = append(files, repositories.File{
				Path:    f.Name,
				Name:    filepath.Base(f.Name),
				Type:    filepath.Ext(f.Name),
				Content: content,
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate files: %w", err)
	}

	return &repositories.Repository{
		Name:  repoDir,
		Files: files,
		SHA:   headRef.Hash().String(),
		URL:   githubURL,
	}, nil
}

func shouldIncludeFile(filePath string, fileTypes []string) bool {
	// If no file types specified, include all files
	if len(fileTypes) == 0 {
		return true
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	// Remove leading dot from extension
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

	for _, fileType := range fileTypes {
		if strings.EqualFold(fileType, ext) {
			return true
		}
	}

	return false
}
