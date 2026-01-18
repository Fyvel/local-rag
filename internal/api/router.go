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
	listDiscussionsUseCase := cfg.Container.NewListDiscussionsUseCase()
	createDiscussionUseCase := cfg.Container.NewCreateDiscussionUseCase()
	getDiscussionUseCase := cfg.Container.NewGetDiscussionUseCase()
	askQuestionUseCase := cfg.Container.NewAskQuestionUseCase()
	askQuestionStreamUseCase := cfg.Container.NewAskQuestionStreamUseCase()
	getDiscussionHistoryUseCase := cfg.Container.NewGetDiscussionHistoryUseCase()
	askQuestionQuickUseCase := cfg.Container.NewAskQuestionQuickUseCase()
	askQuestionQuickStreamUseCase := cfg.Container.NewAskQuestionQuickStreamUseCase()

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", HealthCheckHandler)

		// Indexing endpoints
		v1.POST("/index-github", IndexGithubHandlerFactory(indexRepoUseCase))

		// Quick question endpoints (auto-creates discussions)
		v1.POST("/questions", AskQuestionQuickHandlerFactory(askQuestionQuickUseCase))
		v1.POST("/questions/stream", AskQuestionQuickStreamHandlerFactory(askQuestionQuickStreamUseCase))

		// Discussion endpoints
		v1.GET("/discussions", ListDiscussionsHandlerFactory(listDiscussionsUseCase))
		v1.POST("/discussions", CreateDiscussionHandlerFactory(createDiscussionUseCase))
		v1.GET("/discussions/:id", GetDiscussionHandlerFactory(getDiscussionUseCase))
		v1.POST("/discussions/:id/question", AskQuestionHandlerFactory(askQuestionUseCase))
		v1.POST("/discussions/:id/question/stream", AskQuestionStreamHandlerFactory(askQuestionStreamUseCase))
		v1.GET("/discussions/:id/history", GetDiscussionHistoryHandlerFactory(getDiscussionHistoryUseCase))

		// WebSocket endpoints (alternative to SSE)
		v1.GET("/ws/questions", AskQuestionQuickWSHandlerFactory(askQuestionQuickStreamUseCase))
		v1.GET("/ws/discussions/:id/question", AskQuestionWSHandlerFactory(askQuestionStreamUseCase))
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
