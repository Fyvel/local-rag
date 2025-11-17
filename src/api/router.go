package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Local RAG API
// @version         1.0
// @description     API for local RAG Indexing
// @host            localhost:8080
// @BasePath        /api/v1
func NewRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", HealthCheckHandler)
		v1.POST("/index-github", IndexGithubHandler)

		// v1.GET("/discussions", GetDiscussionsHandler)
		// v1.POST("/discussions", CreateDiscussionHandler)
		// v1.GET("/discussions/:id", GetDiscussionByIDHandler)
		// v1.POST("/discussions/:id/question", CreateQuestionHandler)
		// v1.GET("/discussions/:id/history", GetDiscussionHistoryHandler)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
