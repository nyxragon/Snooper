package pipeline

import (
	"strings"
	"sync"
	"time"

	"github.com/abhi-ingle/cloudsnoop/internal/extractor"
	"github.com/abhi-ingle/cloudsnoop/internal/fetcher"
	"github.com/abhi-ingle/cloudsnoop/internal/matcher"
)

// LinkWithStatus holds a URL and its HTTP status when --check-access is used
type LinkWithStatus struct {
	URL    string
	Status int
}

// Result holds the extraction result for a single URL
type Result struct {
	URL             string
	Links           map[string][]string
	LinksWithStatus map[string][]LinkWithStatus
	Error           error
}

// Pipeline orchestrates fetch -> extract -> match
type Pipeline struct {
	fetcher     *fetcher.Fetcher
	extractor   *extractor.Registry
	matcher     *matcher.Matcher
	workers     int
	delay       int
	checkAccess bool
}

// Config holds pipeline configuration
type Config struct {
	Workers     int
	Timeout     int    // seconds
	Retries     int
	Services    []string
	Delay       int    // milliseconds between requests per worker
	UserAgent   string
	CheckAccess bool
}

// New creates a Pipeline with the given config
func New(cfg Config) (*Pipeline, error) {
	m, err := matcher.New(cfg.Services)
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(cfg.Timeout) * time.Second
	f := fetcher.New(timeout, cfg.Retries, cfg.UserAgent)

	return &Pipeline{
		fetcher:     f,
		extractor:   extractor.NewRegistry(),
		matcher:    m,
		workers:    cfg.Workers,
		delay:      cfg.Delay,
		checkAccess: cfg.CheckAccess,
	}, nil
}

// Run processes URLs concurrently and returns results
func (p *Pipeline) Run(urls []string) []Result {
	// Filter and trim URLs
	var filtered []string
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u != "" {
			filtered = append(filtered, u)
		}
	}
	urls = filtered

	results := make([]Result, len(urls))
	var wg sync.WaitGroup
	sem := make(chan struct{}, p.workers)

	for i, url := range urls {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, u string) {
			defer wg.Done()
			defer func() { <-sem }()
			results[idx] = p.process(u)
		}(i, url)
	}

	wg.Wait()
	p.fetcher.Cleanup()
	return results
}

func (p *Pipeline) process(url string) Result {
	fetched, err := p.fetcher.Fetch(url)
	if err != nil {
		return Result{URL: url, Error: err}
	}

	text, err := p.extractor.Extract(fetched)
	if err != nil {
		return Result{URL: url, Error: err}
	}

	links := p.matcher.Match(text)

	if p.checkAccess && len(links) > 0 {
		linksWithStatus := make(map[string][]LinkWithStatus)
		for svc, urls := range links {
			for _, u := range urls {
				status := p.fetcher.CheckAccess(u)
				linksWithStatus[svc] = append(linksWithStatus[svc], LinkWithStatus{URL: u, Status: status})
			}
		}
		if p.delay > 0 {
			time.Sleep(time.Duration(p.delay) * time.Millisecond)
		}
		return Result{URL: url, LinksWithStatus: linksWithStatus}
	}

	if p.delay > 0 {
		time.Sleep(time.Duration(p.delay) * time.Millisecond)
	}
	return Result{URL: url, Links: links}
}
