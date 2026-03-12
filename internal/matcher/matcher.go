package matcher

import (
	"net/url"
	"regexp"
	"strings"
	"sync"
)

// Matcher extracts cloud storage links from text
type Matcher struct {
	patterns map[string][]*regexp.Regexp
	mu       sync.RWMutex
}

// New creates a Matcher with pre-compiled patterns for the given services
func New(services []string) (*Matcher, error) {
	m := &Matcher{
		patterns: make(map[string][]*regexp.Regexp),
	}

	serviceSet := make(map[string]bool)
	for _, s := range services {
		serviceSet[strings.TrimSpace(s)] = true
	}
	all := serviceSet["all"]
	if all {
		serviceSet = map[string]bool{
			ServiceDrive: true, ServiceSharePoint: true, ServiceDropbox: true,
			ServiceOneDrive: true, ServiceBox: true, ServiceICloud: true,
		}
	}

	for _, p := range AllPatterns() {
		if !serviceSet[p.Service] {
			continue
		}
		re, err := regexp.Compile(p.Regex)
		if err != nil {
			return nil, err
		}
		m.patterns[p.Service] = append(m.patterns[p.Service], re)
	}

	return m, nil
}

// Match extracts and deduplicates cloud links from text
func (m *Matcher) Match(text string) map[string][]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string][]string)
	seen := make(map[string]map[string]bool)

	for service, regexes := range m.patterns {
		if seen[service] == nil {
			seen[service] = make(map[string]bool)
		}
		for _, re := range regexes {
			for _, link := range re.FindAllString(text, -1) {
				normalized := normalizeURL(link)
				if !seen[service][normalized] {
					seen[service][normalized] = true
					result[service] = append(result[service], link)
				}
			}
		}
	}

	return result
}

// normalizeURL normalizes a URL for deduplication (e.g. strip query params that don't affect content)
func normalizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	// Remove common non-essential query params for deduplication
	q := u.Query()
	q.Del("dl")   // Dropbox dl=0 vs dl=1
	q.Del("raw")  // Dropbox raw=1
	q.Del("usp")  // Google Drive
	u.RawQuery = q.Encode()
	return u.String()
}
