package fetcher

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// InMemoryThreshold is the max size (bytes) to keep in memory
	InMemoryThreshold = 5 * 1024 * 1024 // 5MB
)

// FetchedResult holds the result of a fetch operation
type FetchedResult struct {
	URL         string
	Content     []byte
	FilePath    string // Set if written to disk
	ContentType string
	Filename    string
	UseInMemory bool
}

// Fetcher downloads files from URLs with retries and connection pooling
type Fetcher struct {
	client    *http.Client
	retries   int
	timeout   time.Duration
	dumpDir   string
}

// New creates a Fetcher with the given configuration
func New(timeout time.Duration, retries int) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		retries: retries,
		timeout: timeout,
	}
}

// Fetch downloads the URL and returns content (in-memory or disk path)
func (f *Fetcher) Fetch(url string) (*FetchedResult, error) {
	var lastErr error
	backoff := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}

	for attempt := 0; attempt <= f.retries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff[attempt-1])
		}

		result, err := f.doFetch(url)
		if err == nil {
			return result, nil
		}
		lastErr = err

		// Retry on transient errors
		if !isRetryable(err) {
			return nil, err
		}
	}

	return nil, lastErr
}

func (f *Fetcher) doFetch(url string) (*FetchedResult, error) {
	var req *http.Request
	var err error
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Snooper/1.0")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPError{StatusCode: resp.StatusCode, URL: url}
	}

	contentType := resp.Header.Get("Content-Type")
	filename := filenameFromURL(url, resp.Header.Get("Content-Disposition"))

	// Read first chunk to check size
	buf := make([]byte, InMemoryThreshold+1)
	n, err := io.ReadFull(resp.Body, buf)
	if err == io.EOF {
		return &FetchedResult{
			URL:         url,
			Content:     buf[:n],
			ContentType: contentType,
			Filename:    filename,
			UseInMemory: true,
		}, nil
	}
	if err == io.ErrUnexpectedEOF {
		// Content fits in buffer
		return &FetchedResult{
			URL:         url,
			Content:     buf[:n],
			ContentType: contentType,
			Filename:    filename,
			UseInMemory: true,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	// Content too large - use disk
	if f.dumpDir == "" {
		dir, err := os.MkdirTemp("", "snooper-*")
		if err != nil {
			return nil, err
		}
		f.dumpDir = dir
	}

	filePath := filepath.Join(f.dumpDir, sanitizeFilename(filename))
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(buf[:n])
	if err != nil {
		file.Close()
		os.Remove(filePath)
		return nil, err
	}

	_, err = io.Copy(file, resp.Body)
	file.Close()
	if err != nil {
		os.Remove(filePath)
		return nil, err
	}

	return &FetchedResult{
		URL:         url,
		FilePath:    filePath,
		ContentType: contentType,
		Filename:    filename,
		UseInMemory: false,
	}, nil
}

func filenameFromURL(urlStr, contentDisposition string) string {
	// Try Content-Disposition first
	if contentDisposition != "" {
		if idx := strings.Index(contentDisposition, "filename="); idx >= 0 {
			name := strings.Trim(contentDisposition[idx+9:], "\" ")
			if name != "" {
				return name
			}
		}
	}
	// Fallback to URL path
	base := filepath.Base(urlStr)
	if base == "" || base == "." || base == "/" {
		return "downloaded_file"
	}
	return base
}

func sanitizeFilename(name string) string {
	// Remove path traversal and invalid chars
	name = filepath.Base(name)
	name = strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
			return '_'
		}
		return r
	}, name)
	if name == "" {
		return "downloaded_file"
	}
	return name
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	if he, ok := err.(*HTTPError); ok {
		return he.StatusCode >= 500 || he.StatusCode == 429
	}
	return true
}

// HTTPError represents an HTTP error response
type HTTPError struct {
	StatusCode int
	URL        string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.URL)
}

// Cleanup removes temporary files
func (f *Fetcher) Cleanup() {
	if f.dumpDir != "" {
		os.RemoveAll(f.dumpDir)
		f.dumpDir = ""
	}
}

// Reader returns a reader for the fetched content
func (r *FetchedResult) Reader() (io.Reader, error) {
	if r.UseInMemory {
		return bytes.NewReader(r.Content), nil
	}
	return os.Open(r.FilePath)
}

// Bytes returns the content as bytes (reads from disk if needed)
func (r *FetchedResult) Bytes() ([]byte, error) {
	if r.UseInMemory {
		return r.Content, nil
	}
	return os.ReadFile(r.FilePath)
}
