package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"indexing/internal/indexer"
)

type IndexGithubRepoRequest struct {
	GithubURL   string   `json:"github_url"`
	FileTypes   []string `json:"file_types"`
	TargetIndex string   `json:"index_source"`
}

type IndexGithubRepoResponse struct {
	Status  string `json:"status"`
	IndexID string `json:"index_id"`
}

func IndexGithubRepo(w http.ResponseWriter, r *http.Request) {
	var req IndexGithubRepoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.GithubURL == "" {
		http.Error(w, "github_url is required", http.StatusBadRequest)
		return
	}

	valid := validateGithubURL(req.GithubURL)
	if !valid {
		http.Error(w, "Invalid github_url format", http.StatusBadRequest)
		return
	}

	repo, err := indexer.IndexGithubRepository(r.Context(), req.GithubURL, req.FileTypes, req.TargetIndex)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to index the repository: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Indexed repository: %s\n", repo.Name)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(IndexGithubRepoResponse{
		Status:  "success",
		IndexID: req.TargetIndex,
	})
}

func validateGithubURL(url string) bool {
	if url == "" {
		return false
	}
	const githubPrefix = "https://github.com/"
	if len(url) < len(githubPrefix) || url[:len(githubPrefix)] != githubPrefix {
		return false
	}
	return true
}
