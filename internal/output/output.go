package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/abhi-ingle/cloudsnoop/internal/matcher"
	"github.com/abhi-ingle/cloudsnoop/internal/pipeline"
)

// ServiceDisplayNames maps service IDs to display names
var ServiceDisplayNames = map[string]string{
	matcher.ServiceDrive:      "Google Drive",
	matcher.ServiceSharePoint: "SharePoint",
	matcher.ServiceDropbox:    "Dropbox",
	matcher.ServiceOneDrive:   "OneDrive",
	matcher.ServiceBox:        "Box",
	matcher.ServiceICloud:     "iCloud",
}

// Format writes results in the specified format
func Format(w io.Writer, results []pipeline.Result, format string) error {
	switch format {
	case "json":
		return formatJSON(w, results)
	default:
		return formatText(w, results)
	}
}

func formatText(w io.Writer, results []pipeline.Result) error {
	for _, r := range results {
		if r.Error != nil {
			fmt.Fprintf(w, "Processing URL: %s\n", r.URL)
			fmt.Fprintf(w, "Failed: %v\n\n", r.Error)
			continue
		}

		fmt.Fprintf(w, "Processing URL: %s\n", r.URL)

		// Order services for consistent output
		order := []string{
			matcher.ServiceDrive, matcher.ServiceSharePoint, matcher.ServiceDropbox,
			matcher.ServiceOneDrive, matcher.ServiceBox, matcher.ServiceICloud,
		}

		for _, svc := range order {
			links, ok := r.Links[svc]
			if !ok || len(links) == 0 {
				continue
			}
			name := ServiceDisplayNames[svc]
			if name == "" {
				name = svc
			}
			fmt.Fprintf(w, "Found %s links:\n", name)
			for _, link := range links {
				fmt.Fprintln(w, link)
			}
		}
		fmt.Fprintln(w)
	}
	return nil
}

func formatJSON(w io.Writer, results []pipeline.Result) error {
	type urlResult struct {
		URL   string            `json:"url"`
		Links map[string][]string `json:"links,omitempty"`
		Error string            `json:"error,omitempty"`
	}

	var out []urlResult
	for _, r := range results {
		ur := urlResult{URL: r.URL}
		if r.Error != nil {
			ur.Error = r.Error.Error()
		} else {
			ur.Links = r.Links
		}
		out = append(out, ur)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// ParseServices parses comma-separated service list
func ParseServices(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{matcher.ServiceDrive}
	}
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{matcher.ServiceDrive}
	}
	return out
}
