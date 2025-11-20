package textsplitter

import "github.com/tmc/langchaingo/textsplitter"

type MarkdownSplitter struct {
	splitter *textsplitter.MarkdownTextSplitter
}

func NewMarkdownSplitter(opts ...textsplitter.Option) *MarkdownSplitter {
	return &MarkdownSplitter{
		splitter: textsplitter.NewMarkdownTextSplitter(opts...),
	}
}

func (m *MarkdownSplitter) SplitText(text string) ([]string, error) {
	return m.splitter.SplitText(text)
}
