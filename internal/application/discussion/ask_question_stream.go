package discussion

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"local-ai/internal/domain/chat"
	"local-ai/internal/domain/discussion"
)

type AskQuestionStreamUseCase struct {
	repo        discussion.Repository
	chatService chat.Service
}

func NewAskQuestionStreamUseCase(repo discussion.Repository, chatService chat.Service) *AskQuestionStreamUseCase {
	return &AskQuestionStreamUseCase{
		repo:        repo,
		chatService: chatService,
	}
}

type StreamResponse struct {
	MessageID string `json:"message_id,omitempty"`
	Token     string `json:"token,omitempty"`
	Done      bool   `json:"done"`
	Error     string `json:"error,omitempty"`
}

func (uc *AskQuestionStreamUseCase) Execute(ctx context.Context, cmd AskQuestionCommand) (<-chan StreamResponse, error) {
	if cmd.DiscussionID == "" {
		return nil, fmt.Errorf("discussion ID cannot be empty")
	}
	if cmd.Question == "" {
		return nil, fmt.Errorf("question cannot be empty")
	}

	// Retrieve the discussion
	disc, err := uc.repo.FindByID(ctx, cmd.DiscussionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve discussion: %w", err)
	}

	// Create user message
	userMessageID := uuid.New().String()
	userMessage, err := discussion.NewMessage(userMessageID, discussion.RoleUser, cmd.Question)
	if err != nil {
		return nil, fmt.Errorf("failed to create user message: %w", err)
	}

	// Add user message to discussion
	if err := disc.AddMessage(userMessage); err != nil {
		return nil, fmt.Errorf("failed to add user message: %w", err)
	}

	// Build chat history for context
	chatMessages := make([]chat.Message, 0, len(disc.Messages))
	for _, msg := range disc.Messages {
		chatMessages = append(chatMessages, chat.Message{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	// Start streaming completion
	chatReq := chat.CompletionRequest{
		Messages: chatMessages,
		Stream:   true,
	}

	chunkChan, err := uc.chatService.CompleteStreamChannel(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to start streaming: %w", err)
	}

	// Create assistant message ID upfront
	assistantMessageID := uuid.New().String()

	// Create output channel
	outputChan := make(chan StreamResponse, 10)

	// Spawn goroutine to process chunks and manage state
	go func() {
		defer close(outputChan)

		var fullAnswer strings.Builder

		for chunk := range chunkChan {
			// Check for errors
			if chunk.Error != nil {
				outputChan <- StreamResponse{
					Done:  true,
					Error: chunk.Error.Error(),
				}
				return
			}

			// Accumulate content
			if chunk.Content != "" {
				fullAnswer.WriteString(chunk.Content)

				// Emit token
				outputChan <- StreamResponse{
					MessageID: assistantMessageID,
					Token:     chunk.Content,
					Done:      false,
				}
			}

			if chunk.Done {
				// Create and save assistant message
				assistantMessage, err := discussion.NewMessage(
					assistantMessageID,
					discussion.RoleAssistant,
					fullAnswer.String(),
				)
				if err != nil {
					outputChan <- StreamResponse{
						Done:  true,
						Error: fmt.Sprintf("failed to create assistant message: %v", err),
					}
					return
				}

				// Add assistant message to discussion
				if err := disc.AddMessage(assistantMessage); err != nil {
					outputChan <- StreamResponse{
						Done:  true,
						Error: fmt.Sprintf("failed to add assistant message: %v", err),
					}
					return
				}

				// Save updated discussion
				if err := uc.repo.Save(ctx, disc); err != nil {
					outputChan <- StreamResponse{
						Done:  true,
						Error: fmt.Sprintf("failed to save discussion: %v", err),
					}
					return
				}

				// Send final done message
				outputChan <- StreamResponse{
					MessageID: assistantMessageID,
					Done:      true,
				}
				return
			}
		}
	}()

	return outputChan, nil
}
