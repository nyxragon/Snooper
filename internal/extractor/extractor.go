package extractor

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"

	"github.com/abhi-ingle/cloudsnoop/internal/fetcher"
)

// Extractor extracts text from fetched content
type Extractor interface {
	Extract(result *fetcher.FetchedResult) (string, error)
}

// Registry maps file types to extractors
type Registry struct {
	byExt     map[string]Extractor
	byMIME    map[string]Extractor
	defaultExtractor Extractor
}

// NewRegistry creates a registry with all built-in extractors
func NewRegistry() *Registry {
	r := &Registry{
		byExt:  make(map[string]Extractor),
		byMIME: make(map[string]Extractor),
	}

	txt := &TextExtractor{}
	doc := &OfficeExtractor{}
	html := &HTMLExtractor{}

	for _, ext := range []string{".txt"} {
		r.byExt[ext] = txt
	}
	for _, ext := range []string{".pdf"} {
		r.byExt[ext] = &PDFExtractor{}
	}
	for _, ext := range []string{".pptx", ".docx", ".xlsx", ".odt", ".ods"} {
		r.byExt[ext] = doc
	}
	for _, ext := range []string{".html", ".htm"} {
		r.byExt[ext] = html
	}

	// MIME type fallbacks
	r.byMIME["text/plain"] = txt
	r.byMIME["application/pdf"] = &PDFExtractor{}
	r.byMIME["application/vnd.openxmlformats-officedocument"] = doc
	r.byMIME["application/vnd.oasis.opendocument"] = doc
	r.byMIME["text/html"] = html
	r.byMIME["application/xhtml"] = html

	r.defaultExtractor = html
	return r
}

// Extract determines the extractor and extracts text
func (r *Registry) Extract(result *fetcher.FetchedResult) (string, error) {
	ext := strings.ToLower(filepath.Ext(result.Filename))
	if ext == "" {
		ext = extFromURL(result.URL)
	}

	if ex, ok := r.byExt[ext]; ok {
		return ex.Extract(result)
	}

	// Try Content-Type
	if result.ContentType != "" {
		mime := strings.Split(result.ContentType, ";")[0]
		mime = strings.TrimSpace(strings.ToLower(mime))
		if ex, ok := r.byMIME[mime]; ok {
			return ex.Extract(result)
		}
		if strings.HasPrefix(mime, "application/vnd.openxmlformats-officedocument") ||
			strings.HasPrefix(mime, "application/vnd.oasis.opendocument") ||
			strings.HasPrefix(mime, "application/vnd.ms-excel") {
			return r.byExt[".docx"].Extract(result)
		}
	}

	return r.defaultExtractor.Extract(result)
}

func extFromURL(url string) string {
	path := url
	if idx := strings.Index(path, "?"); idx >= 0 {
		path = path[:idx]
	}
	return strings.ToLower(filepath.Ext(path))
}

// readResult reads content from FetchedResult into bytes
func readResult(result *fetcher.FetchedResult) ([]byte, error) {
	if result.UseInMemory {
		return result.Content, nil
	}
	return result.Bytes()
}

// readerFromResult returns an io.Reader for the result
func readerFromResult(result *fetcher.FetchedResult) (io.Reader, error) {
	b, err := readResult(result)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
