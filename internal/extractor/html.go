package extractor

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/abhi-ingle/cloudsnoop/internal/fetcher"
)

// HTMLExtractor extracts text and href links from HTML
type HTMLExtractor struct{}

// Extract extracts body text and all href attributes (for link discovery)
func (e *HTMLExtractor) Extract(result *fetcher.FetchedResult) (string, error) {
	b, err := readResult(result)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	// Body text
	sb.WriteString(doc.Find("body").Text())
	sb.WriteString("\n")

	// Extract href attributes (often contain cloud links)
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok && href != "" {
			sb.WriteString(href)
			sb.WriteString("\n")
		}
	})

	return sb.String(), nil
}
