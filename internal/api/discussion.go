package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"local-ai/internal/application/discussion"
	domaindiscussion "local-ai/internal/domain/discussion"
)

// ListDiscussionsResponse is the response structure for listing discussions.
type ListDiscussionsResponse struct {
	Discussions []discussion.DiscussionSummaryDTO `json:"discussions"`
}

// ListDiscussionsHandler handles the GET /discussions endpoint.
// @Summary      List all discussions
// @Description  Retrieves a list of all discussions with their summaries
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Success      200  {object}  SuccessResponse[ListDiscussionsResponse]
// @Failure      500  {object}  ErrorResponse
// @Router       /discussions [get]
func ListDiscussionsHandlerFactory(useCase *discussion.ListDiscussionsUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := useCase.Execute(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to list discussions: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse[ListDiscussionsResponse]{
			Success: true,
			Data: ListDiscussionsResponse{
				Discussions: result.Discussions,
			},
		})
	}
}

// CreateDiscussionRequest is the request structure for creating a discussion.
type CreateDiscussionRequest struct {
	Title string `json:"title" binding:"required" example:"How to use embeddings?"`
}

// CreateDiscussionHandler handles the POST /discussions endpoint.
// @Summary      Create a new discussion
// @Description  Creates a new discussion with the specified title
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Param        request  body      CreateDiscussionRequest  true  "Discussion creation parameters"
// @Success      201      {object}  SuccessResponse[discussion.CreateDiscussionResult]
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /discussions [post]
func CreateDiscussionHandlerFactory(useCase *discussion.CreateDiscussionUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateDiscussionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Invalid request body: %v", err),
			})
			return
		}

		cmd := discussion.CreateDiscussionCommand{
			Title: req.Title,
		}

		result, err := useCase.Execute(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to create discussion: %v", err),
			})
			return
		}

		c.JSON(http.StatusCreated, SuccessResponse[discussion.CreateDiscussionResult]{
			Success: true,
			Data:    *result,
		})
	}
}

// GetDiscussionHandler handles the GET /discussions/:id endpoint.
// @Summary      Get a discussion by ID
// @Description  Retrieves a complete discussion including all messages
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Discussion ID"
// @Success      200  {object}  SuccessResponse[discussion.DiscussionDTO]
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /discussions/{id} [get]
func GetDiscussionHandlerFactory(useCase *discussion.GetDiscussionUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "Discussion ID is required",
			})
			return
		}

		query := discussion.GetDiscussionQuery{
			ID: id,
		}

		result, err := useCase.Execute(c.Request.Context(), query)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Discussion not found: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse[discussion.DiscussionDTO]{
			Success: true,
			Data:    *result,
		})
	}
}

// AskQuestionRequest is the request structure for asking a question.
type AskQuestionRequest struct {
	Question string `json:"question" binding:"required" example:"What are vector embeddings?"`
}

// AskQuestionHandler handles the POST /discussions/:id/question endpoint.
// @Summary      Ask a question in a discussion
// @Description  Adds a user question to a discussion
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Discussion ID"
// @Param        request  body      AskQuestionRequest    true  "Question parameters"
// @Success      201      {object}  SuccessResponse[discussion.AskQuestionResult]
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /discussions/{id}/question [post]
func AskQuestionHandlerFactory(useCase *discussion.AskQuestionUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "Discussion ID is required",
			})
			return
		}

		var req AskQuestionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Invalid request body: %v", err),
			})
			return
		}

		cmd := discussion.AskQuestionCommand{
			DiscussionID: id,
			Question:     req.Question,
		}

		result, err := useCase.Execute(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to add question: %v", err),
			})
			return
		}

		c.JSON(http.StatusCreated, SuccessResponse[discussion.AskQuestionResult]{
			Success: true,
			Data:    *result,
		})
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
// @Success      200      {string}  string                "Streaming response"
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /discussions/{id}/question/stream [post]
func AskQuestionStreamHandlerFactory(useCase *discussion.AskQuestionUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "Discussion ID is required",
			})
			return
		}

		var req AskQuestionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Invalid request body: %v", err),
			})
			return
		}

		cmd := discussion.AskQuestionCommand{
			DiscussionID: id,
			Question:     req.Question,
		}

		stream, disc, assistantMessageID, err := useCase.ExecuteStream(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to ask question: %v", err),
			})
			return
		}
		defer stream.Close()

		// Set SSE headers
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Transfer-Encoding", "chunked")

		// Buffer to collect the full response for saving
		var fullResponse string

		// Stream the response
		buf := make([]byte, 1024)
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   "Streaming not supported",
			})
			return
		}

		for {
			n, err := stream.Read(buf)
			if n > 0 {
				chunk := string(buf[:n])
				fullResponse += chunk

				// Send SSE event
				fmt.Fprintf(c.Writer, "data: %s\n\n", chunk)
				flusher.Flush()
			}

			if err != nil {
				if err.Error() != "EOF" {
					fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", err.Error())
					flusher.Flush()
				}
				break
			}
		}

		// Save the complete assistant message to the discussion
		if fullResponse != "" {
			assistantMessage, err := domaindiscussion.NewMessage(assistantMessageID, domaindiscussion.RoleAssistant, fullResponse)
			if err == nil {
				disc.AddMessage(assistantMessage)
				useCase.SaveDiscussion(c.Request.Context(), disc)
			}
		}

		// Send completion event
		fmt.Fprintf(c.Writer, "event: done\ndata: {\"message_id\": \"%s\"}\n\n", assistantMessageID)
		flusher.Flush()
	}
}

// GetDiscussionHistoryHandler handles the GET /discussions/:id/history endpoint.
// @Summary      Get discussion message history
// @Description  Retrieves all messages in a discussion
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Discussion ID"
// @Success      200  {object}  SuccessResponse[discussion.GetDiscussionHistoryResult]
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /discussions/{id}/history [get]
func GetDiscussionHistoryHandlerFactory(useCase *discussion.GetDiscussionHistoryUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "Discussion ID is required",
			})
			return
		}

		query := discussion.GetDiscussionHistoryQuery{
			DiscussionID: id,
		}

		result, err := useCase.Execute(c.Request.Context(), query)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Discussion not found: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse[discussion.GetDiscussionHistoryResult]{
			Success: true,
			Data:    *result,
		})
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
// @Success      201      {object}  SuccessResponse[discussion.AskQuestionQuickResult]
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /questions [post]
func AskQuestionQuickHandlerFactory(useCase *discussion.AskQuestionQuickUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req QuickQuestionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Invalid request body: %v", err),
			})
			return
		}

		cmd := discussion.AskQuestionQuickCommand{
			DiscussionID: req.DiscussionID,
			Question:     req.Question,
		}

		result, err := useCase.Execute(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to ask question: %v", err),
			})
			return
		}

		c.JSON(http.StatusCreated, SuccessResponse[discussion.AskQuestionQuickResult]{
			Success: true,
			Data:    *result,
		})
	}
}

// AskQuestionQuickStreamHandler handles the POST /questions/stream endpoint with SSE.
// @Summary      Ask a question with streaming (with auto-discussion creation)
// @Description  Asks a question with streaming response, creating a new discussion if discussion_id is not provided
// @Tags         discussions
// @Accept       json
// @Produce      text/event-stream
// @Param        request  body      QuickQuestionRequest  true  "Question parameters"
// @Success      200      {string}  string                "Streaming response"
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /questions/stream [post]
func AskQuestionQuickStreamHandlerFactory(useCase *discussion.AskQuestionQuickUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req QuickQuestionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Invalid request body: %v", err),
			})
			return
		}

		cmd := discussion.AskQuestionQuickCommand{
			DiscussionID: req.DiscussionID,
			Question:     req.Question,
		}

		stream, disc, assistantMessageID, isNew, err := useCase.ExecuteStream(c.Request.Context(), cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to ask question: %v", err),
			})
			return
		}
		defer stream.Close()

		// Set SSE headers
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Transfer-Encoding", "chunked")

		// Send initial metadata about the discussion
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   "Streaming not supported",
			})
			return
		}

		// Send discussion info
		fmt.Fprintf(c.Writer, "event: discussion\ndata: {\"discussion_id\": \"%s\", \"title\": \"%s\", \"is_new\": %t}\n\n", disc.ID, disc.Title, isNew)
		flusher.Flush()

		// Buffer to collect the full response for saving
		var fullResponse string

		// Stream the response
		buf := make([]byte, 1024)
		for {
			n, err := stream.Read(buf)
			if n > 0 {
				chunk := string(buf[:n])
				fullResponse += chunk

				// Send SSE event
				fmt.Fprintf(c.Writer, "data: %s\n\n", chunk)
				flusher.Flush()
			}

			if err != nil {
				if err.Error() != "EOF" {
					fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", err.Error())
					flusher.Flush()
				}
				break
			}
		}

		// Save the complete assistant message to the discussion
		if fullResponse != "" {
			assistantMessage, err := domaindiscussion.NewMessage(assistantMessageID, domaindiscussion.RoleAssistant, fullResponse)
			if err == nil {
				disc.AddMessage(assistantMessage)
				useCase.SaveDiscussion(c.Request.Context(), disc)
			}
		}

		// Send completion event
		fmt.Fprintf(c.Writer, "event: done\ndata: {\"message_id\": \"%s\", \"discussion_id\": \"%s\"}\n\n", assistantMessageID, disc.ID)
		flusher.Flush()
	}
}
