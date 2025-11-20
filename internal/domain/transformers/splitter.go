package transformers

// TextChunker is a domain interface for splitting text into chunks.
// This abstraction allows the application layer to remain independent
// of specific text splitting implementations.
type TextChunker interface {
	// SplitText splits the input text into chunks according to
	// implementation-specific rules (e.g., by token count, character count,
	// markdown structure, etc.)
	SplitText(text string) ([]string, error)
}
