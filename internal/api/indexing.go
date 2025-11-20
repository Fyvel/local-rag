package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"local-ai/internal/application/indexing"
)

type IndexGithubRepoRequest struct {
	GithubURL   string   `json:"github_url" binding:"required" example:"https://github.com/username/repo"`
	FileTypes   []string `json:"file_types" example:"md,mdx"`
	TargetIndex string   `json:"target_index" example:"target_index"`
}

type IndexGithubRepoResponse struct {
	IndexID string `json:"index_id" example:"target_index"`
}

// IndexGithubHandler handles the GitHub repository indexing endpoint.
// @Summary      Index a GitHub repository
// @Description  Fetches and indexes the contents of a GitHub repository
// @Tags         indexing
// @Accept       json
// @Produce      json
// @Param        request  body      IndexGithubRepoRequest  true  "Repository indexing parameters"
// @Success      201      {object}  SuccessResponse[IndexGithubRepoResponse]
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /index-github [post]
func IndexGithubHandlerFactory(useCase *indexing.IndexRepositoryUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req IndexGithubRepoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Invalid request body: %v", err),
			})
			return
		}

		// Validate GitHub URL format
		if !validateGithubURL(req.GithubURL) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "Invalid github_url format. Must be a valid GitHub repository URL (https://github.com/owner/repo)",
			})
			return
		}

		// Execute use case
		cmd := indexing.IndexRepositoryCommand{
			GithubURL:   req.GithubURL,
			FileTypes:   req.FileTypes,
			TargetIndex: req.TargetIndex,
		}

		result, err := useCase.Execute(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to index repository: %v", err),
			})
			return
		}

		c.JSON(http.StatusCreated, SuccessResponse[IndexGithubRepoResponse]{
			Success: true,
			Data: IndexGithubRepoResponse{
				IndexID: result.IndexID,
			},
		})
	}
}

func validateGithubURL(url string) bool {
	if url == "" {
		return false
	}

	const githubPrefix = "https://github.com/"
	if !strings.HasPrefix(url, githubPrefix) {
		return false
	}

	// Check that there's content after the prefix (owner/repo)
	remaining := strings.TrimPrefix(url, githubPrefix)
	parts := strings.Split(remaining, "/")
	return len(parts) >= 2 && parts[0] != "" && parts[1] != ""
}
