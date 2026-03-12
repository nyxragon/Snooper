package extractor

import (
	"os"
	"path/filepath"
	"strings"

	"code.sajari.com/docconv/v2"
	"github.com/abhi-ingle/cloudsnoop/internal/fetcher"
)

// OfficeExtractor extracts text from Office documents (pptx, docx, xlsx, odt, ods)
type OfficeExtractor struct{}

var officeExts = map[string]bool{
	".pptx": true, ".docx": true, ".xlsx": true,
	".odt": true, ".ods": true,
}

// Extract uses docconv for Office formats
func (e *OfficeExtractor) Extract(result *fetcher.FetchedResult) (string, error) {
	ext := strings.ToLower(filepath.Ext(result.Filename))
	if ext == "" {
		ext = extFromURL(result.URL)
	}
	if !officeExts[ext] {
		ext = ".docx" // docconv default
	}

	var path string
	if result.UseInMemory {
		tmp, err := os.CreateTemp("", "snooper-*"+ext)
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
