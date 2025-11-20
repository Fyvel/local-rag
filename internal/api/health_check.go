package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheckResponse struct {
	Status  string `json:"status" example:"ok"`
	Service string `json:"service" example:"local-rag-api"`
	Version string `json:"version" example:"1.0.0"`
}

// HealthCheckHandler handles the health check endpoint.
// @Summary      Health check endpoint
// @Description  Returns the health status of the API service
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  SuccessResponse[HealthCheckResponse]
// @Failure      500  {object}  ErrorResponse
// @Router       /health [get]
func HealthCheckHandler(c *gin.Context) {
	response := HealthCheckResponse{
		Status:  "ok",
		Service: "local-rag-api",
		Version: "1.0.0",
	}

	c.JSON(http.StatusOK, SuccessResponse[HealthCheckResponse]{
		Success: true,
		Data:    response,
	})
}
