package discussion

// =============================================================================
// Shared Result DTOs
// =============================================================================

// MessageResult is a data transfer object for messages.
// Used by GetDiscussionUseCase and GetDiscussionHistoryUseCase.
type MessageResult struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// DiscussionResult is a data transfer object for a complete discussion.
// Used by GetDiscussionUseCase.
type DiscussionResult struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Messages  []MessageResult `json:"messages"`
	Status    string          `json:"status"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

// DiscussionSummaryResult is a data transfer object for discussion summaries.
// Used by ListDiscussionsUseCase.
type DiscussionSummaryResult struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	MessageCount int    `json:"message_count"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// =============================================================================
// Command DTOs
// =============================================================================

// CreateDiscussionCommand contains the parameters for creating a discussion.
type CreateDiscussionCommand struct {
	Title string
}

// AskQuestionCommand contains the parameters for asking a question.
type AskQuestionCommand struct {
	DiscussionID string
	Question     string
}

// AskQuestionQuickCommand contains the parameters for a quick question.
type AskQuestionQuickCommand struct {
	DiscussionID string // Optional - if empty, creates a new discussion
	Question     string
}

// =============================================================================
// Query DTOs
// =============================================================================

// GetDiscussionQuery contains the parameters for getting a discussion.
type GetDiscussionQuery struct {
	ID string
}

// GetDiscussionHistoryQuery contains the parameters for getting discussion history.
type GetDiscussionHistoryQuery struct {
	DiscussionID string
}

// =============================================================================
// Use Case Result DTOs
// =============================================================================

// CreateDiscussionResult contains the result of creating a discussion.
type CreateDiscussionResult struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
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

// AskQuestionQuickResult contains the result of a quick question.
type AskQuestionQuickResult struct {
	DiscussionID       string `json:"discussion_id"`
	DiscussionTitle    string `json:"discussion_title"`
	UserMessageID      string `json:"user_message_id"`
	AssistantMessageID string `json:"assistant_message_id"`
	Question           string `json:"question"`
	Answer             string `json:"answer"`
	QuestionTimestamp  string `json:"question_timestamp"`
	AnswerTimestamp    string `json:"answer_timestamp"`
	IsNewDiscussion    bool   `json:"is_new_discussion"`
}

// GetDiscussionHistoryResult contains the message history of a discussion.
type GetDiscussionHistoryResult struct {
	DiscussionID string          `json:"discussion_id"`
	Messages     []MessageResult `json:"messages"`
}

// =============================================================================
// Streaming Response DTOs
// =============================================================================

// StreamResponse is the streaming response for AskQuestionStreamUseCase.
type StreamResponse struct {
	MessageID string `json:"message_id,omitempty"`
	Token     string `json:"token,omitempty"`
	Done      bool   `json:"done"`
	Error     string `json:"error,omitempty"`
}

// QuickStreamResponse is the streaming response for AskQuestionQuickStreamUseCase.
type QuickStreamResponse struct {
	DiscussionID    string `json:"discussion_id,omitempty"`
	DiscussionTitle string `json:"discussion_title,omitempty"`
	MessageID       string `json:"message_id,omitempty"`
	Token           string `json:"token,omitempty"`
	Done            bool   `json:"done"`
	IsNewDiscussion bool   `json:"is_new_discussion,omitempty"`
	Error           string `json:"error,omitempty"`
}
