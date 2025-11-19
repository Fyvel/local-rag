package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"local-ai/internal/application/indexing"
	"local-ai/internal/di"
)

type IndexGithubRepoRequest struct {
	GithubURL   string   `json:"github_url" binding:"required" example:"https://github.com/username/repo"`
	FileTypes   []string `json:"file_types" example:"md,mdx"`
	TargetIndex string   `json:"index_source" example:"target_index"`
}

type IndexGithubRepoResponse struct {
	IndexID string `json:"index_id" example:"target_index"`
}

// @Summary      Index a GitHub repository
// @Description  Fetches and indexes the contents of a GitHub repository
// @Tags         indexing
// @Accept       json
// @Produce      json
// @Param        request  body      IndexGithubRepoRequest  true  "Repository indexing parameters"
// @Success      201      {object}  SuccessResponse[IndexGithubRepoResponse]
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /api/v1/index-github [post]
func IndexGithubHandler(c *gin.Context) {
	var req IndexGithubRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if req.GithubURL == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "github_url is required",
		})
		return
	}

	if !validateGithubURL(req.GithubURL) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid github_url format",
		})
		return
	}

	// TODO: Inject container via dependency injection instead of creating here
	// For now, create container with default config
	container := createContainer(req.TargetIndex)

	// Create use case with dependencies
	useCase := container.NewIndexRepositoryUseCase()

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

func createContainer(targetIndex string) *di.Container {
	cfg := di.Config{
		ChromaURL:   "http://localhost:8000",
		OllamaURL:   "http://localhost:11434/api",
		TargetIndex: targetIndex,
	}
	return di.NewContainer(cfg)
}
