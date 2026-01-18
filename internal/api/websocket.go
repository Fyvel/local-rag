package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"local-ai/internal/application/discussion"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type WSMessage struct {
	Type  string                          `json:"type"` // "token", "done", "error", "discussion"
	Data  interface{}                     `json:"data,omitempty"`
	Token string                          `json:"token,omitempty"`
	Error string                          `json:"error,omitempty"`
	Meta  *discussion.QuickStreamResponse `json:"meta,omitempty"`
}

// @Summary      Ask a question with WebSocket streaming
// @Description  Asks a question with WebSocket streaming, creating a new discussion if needed
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Param        discussion_id  query     string  false  "Discussion ID (optional)"
// @Param        question       query     string  true   "Question to ask"
// @Success      101            {string}  string  "Switching Protocols"
// @Router       /ws/questions [get]
func AskQuestionQuickWSHandlerFactory(useCase *discussion.AskQuestionQuickStreamUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Upgrade HTTP connection to WebSocket
		conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("Failed to upgrade to WebSocket: %v", err)))
			return
		}
		defer conn.Close()

		// Read initial message with question
		var req QuickQuestionRequest
		if err := conn.ReadJSON(&req); err != nil {
			sendWSError(conn, fmt.Sprintf("Failed to read question: %v", err))
			return
		}

		if req.Question == "" {
			sendWSError(conn, "Question cannot be empty")
			return
		}

		cmd := discussion.AskQuestionQuickCommand{
			DiscussionID: req.DiscussionID,
			Question:     req.Question,
		}

		// Start streaming
		responseChan, err := useCase.Execute(c.Request.Context(), cmd)
		if err != nil {
			sendWSError(conn, fmt.Sprintf("Failed to start streaming: %v", err))
			return
		}

		// Stream responses
		for response := range responseChan {
			if response.Error != "" {
				sendWSError(conn, response.Error)
				return
			}

			// Send discussion metadata if present
			if response.IsNewDiscussion && response.DiscussionID != "" && response.Token == "" {
				msg := WSMessage{
					Type: "discussion",
					Data: map[string]interface{}{
						"discussion_id": response.DiscussionID,
						"title":         response.DiscussionTitle,
						"is_new":        true,
					},
				}
				if err := conn.WriteJSON(msg); err != nil {
					return
				}
				continue
			}

			if response.Done {
				msg := WSMessage{
					Type: "done",
					Data: map[string]interface{}{
						"message_id":    response.MessageID,
						"discussion_id": response.DiscussionID,
					},
				}
				if err := conn.WriteJSON(msg); err != nil {
					return
				}
				return
			}

			if response.Token != "" {
				msg := WSMessage{
					Type:  "token",
					Token: response.Token,
				}
				if err := conn.WriteJSON(msg); err != nil {
					return
				}
			}
		}
	}
}

// AskQuestionWSHandler handles WebSocket streaming for discussion questions.
// @Summary      Ask a question in a discussion with WebSocket streaming
// @Description  Asks a question in an existing discussion with WebSocket streaming
// @Tags         discussions
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Discussion ID"
// @Success      101 {string}  string  "Switching Protocols"
// @Router       /ws/discussions/{id}/question [get]
func AskQuestionWSHandlerFactory(useCase *discussion.AskQuestionStreamUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		discussionID := c.Param("id")
		if discussionID == "" {
			c.JSON(http.StatusBadRequest, NewErrorResponse("Discussion ID is required"))
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("Failed to upgrade to WebSocket: %v", err)))
			return
		}
		defer conn.Close()

		// Read initial message with question
		var req AskQuestionRequest
		if err := conn.ReadJSON(&req); err != nil {
			sendWSError(conn, fmt.Sprintf("Failed to read question: %v", err))
			return
		}

		if req.Question == "" {
			sendWSError(conn, "Question cannot be empty")
			return
		}

		cmd := discussion.AskQuestionCommand{
			DiscussionID: discussionID,
			Question:     req.Question,
		}

		// Start streaming
		responseChan, err := useCase.Execute(c.Request.Context(), cmd)
		if err != nil {
			sendWSError(conn, fmt.Sprintf("Failed to start streaming: %v", err))
			return
		}

		// Stream responses
		for response := range responseChan {
			if response.Error != "" {
				sendWSError(conn, response.Error)
				return
			}

			if response.Done {
				msg := WSMessage{
					Type: "done",
					Data: map[string]interface{}{
						"message_id": response.MessageID,
					},
				}
				if err := conn.WriteJSON(msg); err != nil {
					return
				}
				return
			}

			if response.Token != "" {
				msg := WSMessage{
					Type:  "token",
					Token: response.Token,
				}
				if err := conn.WriteJSON(msg); err != nil {
					return
				}
			}
		}
	}
}

func sendWSError(conn *websocket.Conn, errMsg string) {
	msg := WSMessage{
		Type:  "error",
		Error: errMsg,
	}
	conn.WriteJSON(msg)
}
