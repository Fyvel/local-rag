package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"local-ai/internal/di"
)

// RouterConfig holds configuration for the API router.
type RouterConfig struct {
	Container *di.Container
}

// @title           Local RAG API
// @version         1.0
// @description     API for local RAG indexing and document management
// @host            localhost:8080
// @BasePath        /api/v1
func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.Default()

	// Create use cases from container (dependency injection)
	indexRepoUseCase := cfg.Container.NewIndexRepositoryUseCase()

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", HealthCheckHandler)

		// Indexing endpoints
		v1.POST("/index-github", IndexGithubHandlerFactory(indexRepoUseCase))

		// Future discussion endpoints (commented out for now)
		// v1.GET("/discussions", GetDiscussionsHandler)
		// v1.POST("/discussions", CreateDiscussionHandler)
		// v1.GET("/discussions/:id", GetDiscussionByIDHandler)
		// v1.POST("/discussions/:id/question", CreateQuestionHandler)
		// v1.GET("/discussions/:id/history", GetDiscussionHistoryHandler)
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
