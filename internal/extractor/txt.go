package extractor

import (
	"github.com/abhi-ingle/cloudsnoop/internal/fetcher"
)

// TextExtractor extracts text from plain text files
type TextExtractor struct{}

// Extract returns the raw content as text
func (e *TextExtractor) Extract(result *fetcher.FetchedResult) (string, error) {
	b, err := readResult(result)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
