package discussion

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"local-ai/internal/domain/chat"
	"local-ai/internal/domain/discussion"
)

type AskQuestionQuickStreamUseCase struct {
	repo        discussion.Repository
	chatService chat.Service
	titleModel  string
}

func NewAskQuestionQuickStreamUseCase(repo discussion.Repository, chatService chat.Service, titleModel string) *AskQuestionQuickStreamUseCase {
	if titleModel == "" {
		titleModel = "mistral:latest"
	}
	return &AskQuestionQuickStreamUseCase{
		repo:        repo,
		chatService: chatService,
		titleModel:  titleModel,
	}
}

type QuickStreamResponse struct {
	DiscussionID    string `json:"discussion_id,omitempty"`
	DiscussionTitle string `json:"discussion_title,omitempty"`
	MessageID       string `json:"message_id,omitempty"`
	Token           string `json:"token,omitempty"`
	Done            bool   `json:"done"`
	IsNewDiscussion bool   `json:"is_new_discussion,omitempty"`
	Error           string `json:"error,omitempty"`
}

func (uc *AskQuestionQuickStreamUseCase) Execute(ctx context.Context, cmd AskQuestionQuickCommand) (<-chan QuickStreamResponse, error) {
	if cmd.Question == "" {
		return nil, fmt.Errorf("question cannot be empty")
	}

	outputChan := make(chan QuickStreamResponse, 10)

	// Spawn goroutine to handle async discussion creation and streaming
	go func() {
		defer close(outputChan)

		var disc *discussion.Discussion
		var isNew bool

		// Check if discussion exists or needs to be created
		if cmd.DiscussionID != "" {
			// Try to retrieve existing discussion
			existingDisc, err := uc.repo.FindByID(ctx, cmd.DiscussionID)
			if err != nil {
				outputChan <- QuickStreamResponse{
					Done:  true,
					Error: fmt.Sprintf("discussion with ID %s not found: %v", cmd.DiscussionID, err),
				}
				return
			}
			disc = existingDisc
			isNew = false
		} else {
			// Create new discussion with AI-generated title
			title, err := uc.chatService.GenerateTitle(ctx, cmd.Question, uc.titleModel)
			if err != nil {
				outputChan <- QuickStreamResponse{
					Done:  true,
					Error: fmt.Sprintf("failed to generate discussion title: %v", err),
				}
				return
			}

			// Create discussion
			newDiscussionID := uuid.New().String()
			newDisc, err := discussion.NewDiscussion(newDiscussionID, title)
			if err != nil {
				outputChan <- QuickStreamResponse{
					Done:  true,
					Error: fmt.Sprintf("failed to create discussion: %v", err),
				}
				return
			}

			disc = newDisc
			isNew = true

			// Emit discussion metadata
			outputChan <- QuickStreamResponse{
				DiscussionID:    disc.ID,
				DiscussionTitle: disc.Title,
				IsNewDiscussion: true,
				Done:            false,
			}
		}

		// Create user message
		userMessageID := uuid.New().String()
		userMessage, err := discussion.NewMessage(userMessageID, discussion.RoleUser, cmd.Question)
		if err != nil {
			outputChan <- QuickStreamResponse{
				Done:  true,
				Error: fmt.Sprintf("failed to create user message: %v", err),
			}
			return
		}

		// Add user message to discussion
		if err := disc.AddMessage(userMessage); err != nil {
			outputChan <- QuickStreamResponse{
				Done:  true,
				Error: fmt.Sprintf("failed to add user message: %v", err),
			}
			return
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
			outputChan <- QuickStreamResponse{
				Done:  true,
				Error: fmt.Sprintf("failed to start streaming: %v", err),
			}
			return
		}

		// Create assistant message ID upfront
		assistantMessageID := uuid.New().String()

		var fullAnswer strings.Builder

		for chunk := range chunkChan {
			// Check for errors
			if chunk.Error != nil {
				outputChan <- QuickStreamResponse{
					Done:  true,
					Error: chunk.Error.Error(),
				}
				return
			}

			// Accumulate content
			if chunk.Content != "" {
				fullAnswer.WriteString(chunk.Content)

				// Emit token
				outputChan <- QuickStreamResponse{
					DiscussionID:    disc.ID,
					DiscussionTitle: disc.Title,
					MessageID:       assistantMessageID,
					Token:           chunk.Content,
					IsNewDiscussion: isNew,
					Done:            false,
				}
			}

			// Check if done
			if chunk.Done {
				// Create and save assistant message
				assistantMessage, err := discussion.NewMessage(
					assistantMessageID,
					discussion.RoleAssistant,
					fullAnswer.String(),
				)
				if err != nil {
					outputChan <- QuickStreamResponse{
						Done:  true,
						Error: fmt.Sprintf("failed to create assistant message: %v", err),
					}
					return
				}

				// Add assistant message to discussion
				if err := disc.AddMessage(assistantMessage); err != nil {
					outputChan <- QuickStreamResponse{
						Done:  true,
						Error: fmt.Sprintf("failed to add assistant message: %v", err),
					}
					return
				}

				// Save updated discussion
				if err := uc.repo.Save(ctx, disc); err != nil {
					outputChan <- QuickStreamResponse{
						Done:  true,
						Error: fmt.Sprintf("failed to save discussion: %v", err),
					}
					return
				}

				// Send final done message
				outputChan <- QuickStreamResponse{
					DiscussionID:    disc.ID,
					DiscussionTitle: disc.Title,
					MessageID:       assistantMessageID,
					IsNewDiscussion: isNew,
					Done:            true,
				}
				return
			}
		}
	}()

	return outputChan, nil
}
