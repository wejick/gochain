package textsplitter

//go:generate moq -out textsplitter_moq.go . TextSplitter

// TextSplitter split text
type TextSplitter interface {
	SplitText(string) []string
}