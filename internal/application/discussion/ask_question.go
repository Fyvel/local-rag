package discussion

import (
	"context"
	"fmt"
	"io"

	"local-ai/internal/domain/chat"
	"local-ai/internal/domain/discussion"

	"github.com/google/uuid"
)

// AskQuestionUseCase handles adding a question (user message) to a discussion.
type AskQuestionUseCase struct {
	repo        discussion.Repository
	chatService chat.Service
}

// NewAskQuestionUseCase creates a new ask question use case.
func NewAskQuestionUseCase(repo discussion.Repository, chatService chat.Service) *AskQuestionUseCase {
	return &AskQuestionUseCase{
		repo:        repo,
		chatService: chatService,
	}
}

// AskQuestionCommand contains the parameters for asking a question.
type AskQuestionCommand struct {
	DiscussionID string
	Question     string
}

// AskQuestionResult contains the result of asking a question.
type AskQuestionResult struct {
	UserMessageID      string `json:"user_message_id"`
	AssistantMessageID string `json:"assistant_message_id"`
	DiscussionID       string `json:"discussion_id"`
	Question           string `json:"question"`
	Answer             string `json:"answer"`
	QuestionTimestamp  string `json:"question_timestamp"`
	AnswerTimestamp    string `json:"answer_timestamp"`
}

// Execute adds a user question to a discussion and generates an AI response.
func (uc *AskQuestionUseCase) Execute(ctx context.Context, cmd AskQuestionCommand) (*AskQuestionResult, error) {
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

	// Persist the updated discussion
	if err := uc.repo.Save(ctx, disc); err != nil {
		return nil, fmt.Errorf("failed to save discussion: %w", err)
	}

	return &AskQuestionResult{
		UserMessageID:      userMessage.ID,
		AssistantMessageID: assistantMessage.ID,
		DiscussionID:       disc.ID,
		Question:           userMessage.Content,
		Answer:             assistantMessage.Content,
		QuestionTimestamp:  userMessage.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		AnswerTimestamp:    assistantMessage.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// ExecuteStream adds a user question and returns a streaming AI response.
func (uc *AskQuestionUseCase) ExecuteStream(ctx context.Context, cmd AskQuestionCommand) (io.ReadCloser, *discussion.Discussion, string, error) {
	if cmd.DiscussionID == "" {
		return nil, nil, "", fmt.Errorf("discussion ID cannot be empty")
	}
	if cmd.Question == "" {
		return nil, nil, "", fmt.Errorf("question cannot be empty")
	}

	// Retrieve the discussion
	disc, err := uc.repo.FindByID(ctx, cmd.DiscussionID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to retrieve discussion: %w", err)
	}

	// Create user message
	userMessageID := uuid.New().String()
	userMessage, err := discussion.NewMessage(userMessageID, discussion.RoleUser, cmd.Question)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create user message: %w", err)
	}

	// Add user message to discussion
	if err := disc.AddMessage(userMessage); err != nil {
		return nil, nil, "", fmt.Errorf("failed to add user message: %w", err)
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
		return nil, nil, "", fmt.Errorf("failed to generate streaming AI response: %w", err)
	}

	// Generate assistant message ID for tracking
	assistantMessageID := uuid.New().String()

	return stream, disc, assistantMessageID, nil
}

func (uc *AskQuestionUseCase) SaveStreamedResponse(ctx context.Context, disc *discussion.Discussion, assistantMessageID, content string) error {
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
