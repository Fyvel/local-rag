package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheckResponse struct {
	Status  string `json:"status" example:"ok"`
	Service string `json:"service" example:"local-rag-api"`
	Version string `json:"version" example:"1.0.0"`
}

// @Summary      Health check endpoint
// @Description  Returns the health status of the API service
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  SuccessResponse[HealthCheckResponse]
// @Error        500  {object}  ErrorResponse
// @Router       /api/v1/health [get]
func HealthCheckHandler(c *gin.Context) {
	response, err := healthCheck(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorResponse{
				Success: false,
				Error:   err.Error(),
			})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse[HealthCheckResponse]{
		Success: true,
		Data:    response,
	})
}

func healthCheck(ctx context.Context) (HealthCheckResponse, error) {
	return HealthCheckResponse{
		Status:  "ok",
		Service: "local-rag-api",
		Version: "1.0.0",
	}, nil
}
