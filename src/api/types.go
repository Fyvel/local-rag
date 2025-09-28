package api

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Invalid request parameters"`
}

type SuccessResponse[T any] struct {
	Success bool `json:"success" example:"true"`
	Data    T    `json:"data"`
}
