package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"local-ai/internal/application/discussion"
)

// ListDiscussionsHandler handles the GET /discussions endpoint.
// @Summary      List all discussions
// @Description  Retrieves a list of all discussions with their summaries
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Success      200  {object}  Response[[]discussion.DiscussionSummaryResult]
// @Failure      500  {object}  Response[any]
// @Router       /discussions [get]
func ListDiscussionsHandlerFactory(useCase *discussion.ListDiscussionsUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := useCase.Execute(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Sprintf("Failed to list discussions: %v", err)))
			return
		}

		c.JSON(http.StatusOK, NewSuccessResponse(result))
	}
}

// CreateDiscussionHandler handles the POST /discussions endpoint.
// @Summary      Create a new discussion
// @Description  Creates a new discussion with the specified title
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Param        request  body      CreateDiscussionRequest  true  "Discussion creation parameters"
// @Success      201      {object}  Response[discussion.CreateDiscussionResult]
// @Failure      400      {object}  Response[any]
// @Failure      500      {object}  Response[any]
// @Router       /discussions [post]
func CreateDiscussionHandlerFactory(useCase *discussion.CreateDiscussionUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateDiscussionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("Invalid request body: %v", err)))
			return
		}

		cmd := discussion.CreateDiscussionCommand{
			Title: req.Title,
		}

		result, err := useCase.Execute(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Sprintf("Failed to create discussion: %v", err)))
			return
		}

		c.JSON(http.StatusCreated, NewSuccessResponse(result))
	}
}

// GetDiscussionHandler handles the GET /discussions/:id endpoint.
// @Summary      Get a discussion by ID
// @Description  Retrieves a complete discussion including all messages
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Discussion ID"
// @Success      200  {object}  Response[discussion.DiscussionResult]
// @Failure      400  {object}  Response[any]
// @Failure      404  {object}  Response[any]
// @Router       /discussions/{id} [get]
func GetDiscussionHandlerFactory(useCase *discussion.GetDiscussionUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, NewErrorResponse("Discussion ID is required"))
			return
		}

		query := discussion.GetDiscussionQuery{
			ID: id,
		}

		result, err := useCase.Execute(c.Request.Context(), query)
		if err != nil {
			c.JSON(http.StatusNotFound, NewErrorResponse(fmt.Sprintf("Discussion not found: %v", err)))
			return
		}

		c.JSON(http.StatusOK, NewSuccessResponse(result))
	}
}

// AskQuestionHandler handles the POST /discussions/:id/question endpoint.
// @Summary      Ask a question in a discussion
// @Description  Adds a user question to a discussion
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Discussion ID"
// @Param        request  body      AskQuestionRequest    true  "Question parameters"
// @Success      201      {object}  Response[discussion.AskQuestionResult]
// @Failure      400      {object}  Response[any]
// @Failure      404      {object}  Response[any]
// @Failure      500      {object}  Response[any]
// @Router       /discussions/{id}/question [post]
func AskQuestionHandlerFactory(useCase *discussion.AskQuestionUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, NewErrorResponse("Discussion ID is required"))
			return
		}

		var req AskQuestionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("Invalid request body: %v", err)))
			return
		}

		cmd := discussion.AskQuestionCommand{
			DiscussionID: id,
			Question:     req.Question,
		}

		result, err := useCase.Execute(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Sprintf("Failed to add question: %v", err)))
			return
		}

		c.JSON(http.StatusCreated, NewSuccessResponse(result))
	}
}

// AskQuestionStreamHandler handles the POST /discussions/:id/question/stream endpoint with SSE.
// @Summary      Ask a question with streaming response
// @Description  Adds a user question to a discussion and streams the AI response using Server-Sent Events
// @Tags         discussions
// @Accept       json
// @Produce      text/event-stream
// @Param        id       path      string                true  "Discussion ID"
// @Param        request  body      AskQuestionRequest    true  "Question parameters"
// @Success      200      {string}  string                "Streaming response (SSE stream of tokens)"
// @Failure      400      {object}  Response[any]
// @Failure      404      {object}  Response[any]
// @Failure      500      {object}  Response[any]
// @Router       /discussions/{id}/question/stream [post]
func AskQuestionStreamHandlerFactory(useCase *discussion.AskQuestionStreamUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, NewErrorResponse("Discussion ID is required"))
			return
		}

		var req AskQuestionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("Invalid request body: %v", err)))
			return
		}

		cmd := discussion.AskQuestionCommand{
			DiscussionID: id,
			Question:     req.Question,
		}

		// Start streaming
		responseChan, err := useCase.Execute(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Sprintf("Failed to start streaming: %v", err)))
			return
		}

		// Set SSE headers
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Transfer-Encoding", "chunked")

		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Streaming not supported"))
			return
		}

		// Stream tokens from channel
		for response := range responseChan {
			if response.Error != "" {
				// Send error event
				fmt.Fprintf(c.Writer, "event: error\ndata: {\"error\": \"%s\"}\n\n", response.Error)
				flusher.Flush()
				return
			}

			if response.Done {
				// Send completion event
				fmt.Fprintf(c.Writer, "event: done\ndata: {\"message_id\": \"%s\"}\n\n", response.MessageID)
				flusher.Flush()
				return
			}

			if response.Token != "" {
				// Send token event
				fmt.Fprintf(c.Writer, "data: %s\n\n", response.Token)
				flusher.Flush()
			}
		}
	}
}

// GetDiscussionHistoryHandler handles the GET /discussions/:id/history endpoint.
// @Summary      Get discussion message history
// @Description  Retrieves all messages in a discussion
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Discussion ID"
// @Success      200  {object}  Response[discussion.GetDiscussionHistoryResult]
// @Failure      400  {object}  Response[any]
// @Failure      404  {object}  Response[any]
// @Router       /discussions/{id}/history [get]
func GetDiscussionHistoryHandlerFactory(useCase *discussion.GetDiscussionHistoryUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, NewErrorResponse("Discussion ID is required"))
			return
		}

		query := discussion.GetDiscussionHistoryQuery{
			DiscussionID: id,
		}

		result, err := useCase.Execute(c.Request.Context(), query)
		if err != nil {
			c.JSON(http.StatusNotFound, NewErrorResponse(fmt.Sprintf("Discussion not found: %v", err)))
			return
		}

		c.JSON(http.StatusOK, NewSuccessResponse(result))
	}
}

type QuickQuestionRequest struct {
	DiscussionID string `json:"discussion_id,omitempty" example:""`
	Question     string `json:"question" binding:"required" example:"What are vector embeddings?"`
}

// AskQuestionQuickHandler handles the POST /questions endpoint.
// @Summary      Ask a question (with auto-discussion creation)
// @Description  Asks a question, creating a new discussion if discussion_id is not provided. Uses mistral:latest to auto-generate discussion titles.
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Param        request  body      QuickQuestionRequest  true  "Question parameters"
// @Success      201      {object}  Response[discussion.AskQuestionQuickResult]
// @Failure      400      {object}  Response[any]
// @Failure      500      {object}  Response[any]
// @Router       /questions [post]
func AskQuestionQuickHandlerFactory(useCase *discussion.AskQuestionQuickUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req QuickQuestionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("Invalid request body: %v", err)))
			return
		}

		cmd := discussion.AskQuestionQuickCommand{
			DiscussionID: req.DiscussionID,
			Question:     req.Question,
		}

		result, err := useCase.Execute(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Sprintf("Failed to ask question: %v", err)))
			return
		}

		c.JSON(http.StatusCreated, NewSuccessResponse(result))
	}
}

// AskQuestionQuickStreamHandler handles the POST /questions/stream endpoint with SSE.
// @Summary      Ask a question with streaming (with auto-discussion creation)
// @Description  Asks a question with streaming response, creating a new discussion if discussion_id is not provided
// @Tags         discussions
// @Accept       json
// @Produce      text/event-stream
// @Param        request  body      QuickQuestionRequest  true  "Question parameters"
// @Success      200      {string}  string                "Streaming response (SSE stream of tokens)"
// @Failure      400      {object}  Response[any]
// @Failure      500      {object}  Response[any]
// @Router       /questions/stream [post]
func AskQuestionQuickStreamHandlerFactory(useCase *discussion.AskQuestionQuickStreamUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req QuickQuestionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("Invalid request body: %v", err)))
			return
		}

		cmd := discussion.AskQuestionQuickCommand{
			DiscussionID: req.DiscussionID,
			Question:     req.Question,
		}

		// Start streaming
		responseChan, err := useCase.Execute(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Sprintf("Failed to start streaming: %v", err)))
			return
		}

		// Set SSE headers
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Transfer-Encoding", "chunked")

		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Streaming not supported"))
			return
		}

		// Stream responses from channel
		for response := range responseChan {
			if response.Error != "" {
				// Send error event
				fmt.Fprintf(c.Writer, "event: error\ndata: {\"error\": \"%s\"}\n\n", response.Error)
				flusher.Flush()
				return
			}

			// Send discussion metadata if present (new discussion)
			if response.IsNewDiscussion && response.DiscussionID != "" && response.Token == "" {
				fmt.Fprintf(c.Writer, "event: discussion\ndata: {\"discussion_id\": \"%s\", \"title\": \"%s\", \"is_new\": true}\n\n",
					response.DiscussionID, response.DiscussionTitle)
				flusher.Flush()
				continue
			}

			if response.Done {
				// Send completion event
				fmt.Fprintf(c.Writer, "event: done\ndata: {\"message_id\": \"%s\", \"discussion_id\": \"%s\"}\n\n",
					response.MessageID, response.DiscussionID)
				flusher.Flush()
				return
			}

			if response.Token != "" {
				// Send token event
				fmt.Fprintf(c.Writer, "data: %s\n\n", response.Token)
				flusher.Flush()
			}
		}
	}
}
