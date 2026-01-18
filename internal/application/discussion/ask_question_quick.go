package discussion

import (
	"context"
	"fmt"
	"io"

	"local-ai/internal/domain/chat"
	"local-ai/internal/domain/discussion"

	"github.com/google/uuid"
)

type AskQuestionQuickUseCase struct {
	repo        discussion.Repository
	chatService chat.Service
	titleModel  string
}

func NewAskQuestionQuickUseCase(repo discussion.Repository, chatService chat.Service, titleModel string) *AskQuestionQuickUseCase {
	if titleModel == "" {
		titleModel = "mistral:latest"
	}
	return &AskQuestionQuickUseCase{
		repo:        repo,
		chatService: chatService,
		titleModel:  titleModel,
	}
}

func (uc *AskQuestionQuickUseCase) Execute(ctx context.Context, cmd AskQuestionQuickCommand) (*AskQuestionQuickResult, error) {
	if cmd.Question == "" {
		return nil, fmt.Errorf("question cannot be empty")
	}

	var disc *discussion.Discussion
	var isNew bool

	// Check if discussion exists or needs to be created
	if cmd.DiscussionID != "" {
		// Try to retrieve existing discussion
		existingDisc, err := uc.repo.FindByID(ctx, cmd.DiscussionID)
		if err != nil {
			return nil, fmt.Errorf("discussion with ID %s not found: %w", cmd.DiscussionID, err)
		}
		disc = existingDisc
		isNew = false
	} else {
		// Create new discussion with AI-generated title
		title, err := uc.chatService.GenerateTitle(ctx, cmd.Question, uc.titleModel)
		if err != nil {
			return nil, fmt.Errorf("failed to generate discussion title: %w", err)
		}

		// Create discussion
		newDiscussionID := uuid.New().String()
		newDisc, err := discussion.NewDiscussion(newDiscussionID, title)
		if err != nil {
			return nil, fmt.Errorf("failed to create discussion: %w", err)
		}

		disc = newDisc
		isNew = true
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

	// Generate AI response
	chatReq := chat.CompletionRequest{
		Messages: chatMessages,
		Stream:   false,
	}

	response, err := uc.chatService.Complete(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}

	// Create assistant message
	assistantMessageID := uuid.New().String()
	assistantMessage, err := discussion.NewMessage(assistantMessageID, discussion.RoleAssistant, response.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to create assistant message: %w", err)
	}

	// Add assistant message to discussion
	if err := disc.AddMessage(assistantMessage); err != nil {
		return nil, fmt.Errorf("failed to add assistant message: %w", err)
	}

	// Persist the discussion
	if err := uc.repo.Save(ctx, disc); err != nil {
		return nil, fmt.Errorf("failed to save discussion: %w", err)
	}

	return &AskQuestionQuickResult{
		DiscussionID:       disc.ID,
		DiscussionTitle:    disc.Title,
		UserMessageID:      userMessage.ID,
		AssistantMessageID: assistantMessage.ID,
		Question:           userMessage.Content,
		Answer:             assistantMessage.Content,
		QuestionTimestamp:  userMessage.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		AnswerTimestamp:    assistantMessage.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		IsNewDiscussion:    isNew,
	}, nil
}

// ExecuteStream handles asking a question with streaming, creating a discussion if needed.
func (uc *AskQuestionQuickUseCase) ExecuteStream(ctx context.Context, cmd AskQuestionQuickCommand) (io.ReadCloser, *discussion.Discussion, string, bool, error) {
	if cmd.Question == "" {
		return nil, nil, "", false, fmt.Errorf("question cannot be empty")
	}

	var disc *discussion.Discussion
	var isNew bool

	// Check if discussion exists or needs to be created
	if cmd.DiscussionID != "" {
		// Try to retrieve existing discussion
		existingDisc, err := uc.repo.FindByID(ctx, cmd.DiscussionID)
		if err != nil {
			return nil, nil, "", false, fmt.Errorf("discussion with ID %s not found: %w", cmd.DiscussionID, err)
		}
		disc = existingDisc
		isNew = false
	} else {
		// Create new discussion with AI-generated title
		title, err := uc.chatService.GenerateTitle(ctx, cmd.Question, uc.titleModel)
		if err != nil {
			return nil, nil, "", false, fmt.Errorf("failed to generate discussion title: %w", err)
		}

		// Create discussion
		newDiscussionID := uuid.New().String()
		newDisc, err := discussion.NewDiscussion(newDiscussionID, title)
		if err != nil {
			return nil, nil, "", false, fmt.Errorf("failed to create discussion: %w", err)
		}

		disc = newDisc
		isNew = true
	}

	// Create user message
	userMessageID := uuid.New().String()
	userMessage, err := discussion.NewMessage(userMessageID, discussion.RoleUser, cmd.Question)
	if err != nil {
		return nil, nil, "", false, fmt.Errorf("failed to create user message: %w", err)
	}

	// Add user message to discussion
	if err := disc.AddMessage(userMessage); err != nil {
		return nil, nil, "", false, fmt.Errorf("failed to add user message: %w", err)
	}

	// Build chat history for context
	chatMessages := make([]chat.Message, 0, len(disc.Messages))
	for _, msg := range disc.Messages {
		chatMessages = append(chatMessages, chat.Message{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	// Generate streaming AI response
	chatReq := chat.CompletionRequest{
		Messages: chatMessages,
		Stream:   true,
	}

	stream, err := uc.chatService.CompleteStream(ctx, chatReq)
	if err != nil {
		return nil, nil, "", false, fmt.Errorf("failed to generate streaming AI response: %w", err)
	}

	// Generate assistant message ID for tracking
	assistantMessageID := uuid.New().String()

	return stream, disc, assistantMessageID, isNew, nil
}

func (uc *AskQuestionQuickUseCase) SaveStreamedResponse(ctx context.Context, disc *discussion.Discussion, assistantMessageID, content string) error {
	if content == "" {
		return fmt.Errorf("cannot save empty response")
	}

	assistantMessage, err := discussion.NewMessage(assistantMessageID, discussion.RoleAssistant, content)
	if err != nil {
		return fmt.Errorf("failed to create assistant message: %w", err)
	}

	if err := disc.AddMessage(assistantMessage); err != nil {
		return fmt.Errorf("failed to add assistant message: %w", err)
	}

	if err := uc.repo.Save(ctx, disc); err != nil {
		return fmt.Errorf("failed to save discussion: %w", err)
	}

	return nil
}
