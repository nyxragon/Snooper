package extractor

import (
	"os"

	"code.sajari.com/docconv/v2"
	"github.com/abhi-ingle/cloudsnoop/internal/fetcher"
)

// PDFExtractor extracts text from PDF files
type PDFExtractor struct{}

// Extract uses docconv to extract text from PDF
func (e *PDFExtractor) Extract(result *fetcher.FetchedResult) (string, error) {
	var path string
	if result.UseInMemory {
		tmp, err := os.CreateTemp("", "snooper-*.pdf")
		if err != nil {
			return "", err
		}
		defer os.Remove(tmp.Name())
		defer tmp.Close()
		if _, err := tmp.Write(result.Content); err != nil {
			return "", err
		}
		path = tmp.Name()
	} else {
		path = result.FilePath
	}

	res, err := docconv.ConvertPath(path)
	if err != nil {
		return "", err
	}
	return res.Body, nil
}
