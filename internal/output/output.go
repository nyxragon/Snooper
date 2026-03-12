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
func Format(w io.Writer, results []pipeline.Result, format string, checkAccess bool) error {
	switch format {
	case "json":
		return formatJSON(w, results, checkAccess)
	default:
		return formatText(w, results, checkAccess)
	}
}

func formatText(w io.Writer, results []pipeline.Result, checkAccess bool) error {
	order := []string{
		matcher.ServiceDrive, matcher.ServiceSharePoint, matcher.ServiceDropbox,
		matcher.ServiceOneDrive, matcher.ServiceBox, matcher.ServiceICloud,
	}

	for _, r := range results {
		if r.Error != nil {
			fmt.Fprintf(w, "Processing URL: %s\n", r.URL)
			fmt.Fprintf(w, "Failed: %v\n\n", r.Error)
			continue
		}

		fmt.Fprintf(w, "Processing URL: %s\n", r.URL)

		if checkAccess && len(r.LinksWithStatus) > 0 {
			for _, svc := range order {
				links, ok := r.LinksWithStatus[svc]
				if !ok || len(links) == 0 {
					continue
				}
				name := ServiceDisplayNames[svc]
				if name == "" {
					name = svc
				}
				fmt.Fprintf(w, "Found %s links:\n", name)
				for _, l := range links {
					fmt.Fprintf(w, "  %s (%d)\n", l.URL, l.Status)
				}
			}
		} else {
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
		}
		fmt.Fprintln(w)
	}
	return nil
}

func formatJSON(w io.Writer, results []pipeline.Result, checkAccess bool) error {
	type linkEntry struct {
		URL    string `json:"url"`
		Status int    `json:"status,omitempty"`
	}

	type urlResult struct {
		URL             string                   `json:"url"`
		Links           map[string][]string      `json:"links,omitempty"`
		LinksWithStatus map[string][]linkEntry   `json:"links_with_status,omitempty"`
		Error           string                  `json:"error,omitempty"`
	}

	var out []urlResult
	for _, r := range results {
		ur := urlResult{URL: r.URL}
		if r.Error != nil {
			ur.Error = r.Error.Error()
		} else if checkAccess && len(r.LinksWithStatus) > 0 {
			ur.LinksWithStatus = make(map[string][]linkEntry)
			for svc, links := range r.LinksWithStatus {
				for _, l := range links {
					ur.LinksWithStatus[svc] = append(ur.LinksWithStatus[svc], linkEntry{URL: l.URL, Status: l.Status})
				}
			}
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
