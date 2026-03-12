package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/abhi-ingle/cloudsnoop/internal/output"
	"github.com/abhi-ingle/cloudsnoop/internal/pipeline"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var (
	urlFlag      string
	fileFlag     string
	snoopFlag    string
	workers      int
	timeout      int
	retries      int
	outputFmt    string
	delay        int
	userAgent    string
	checkAccess  bool
)

const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/115.0"

var rootCmd = &cobra.Command{
	Use:   "snooper",
	Short: "Extract cloud storage links from URLs and files",
	Long: `Snooper extracts cloud storage links (Google Drive, SharePoint, Dropbox, OneDrive, Box, iCloud)
from web pages, PDFs, Office documents, and text files.`,
	RunE: run,
}

func init() {
	rootCmd.Flags().StringVarP(&urlFlag, "url", "u", "", "Comma-separated URLs to fetch")
	rootCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Path to a file containing URLs (one per line). Use - to read from stdin")
	rootCmd.Flags().StringVarP(&snoopFlag, "snoop", "s", "drive", "Types of links to extract (drive, sharepoint, dropbox, onedrive, box, icloud, all)")
	rootCmd.Flags().IntVarP(&workers, "workers", "w", 5, "Number of concurrent workers")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 60, "HTTP read timeout in seconds")
	rootCmd.Flags().IntVarP(&retries, "retries", "r", 3, "HTTP retry attempts")
	rootCmd.Flags().StringVarP(&outputFmt, "output", "o", "text", "Output format: text or json")
	rootCmd.Flags().IntVarP(&delay, "delay", "d", 0, "Delay in milliseconds between requests per worker (0 = no delay)")
	rootCmd.Flags().StringVarP(&userAgent, "user-agent", "a", defaultUserAgent, "HTTP User-Agent header")
	rootCmd.Flags().BoolVarP(&checkAccess, "check-access", "c", false, "HEAD each discovered link to check accessibility")
}

func run(cmd *cobra.Command, args []string) error {
	var urls []string

	if fileFlag == "-" {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		urls = strings.Split(string(content), "\n")
	} else if fileFlag != "" {
		content, err := os.ReadFile(fileFlag)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		urls = strings.Split(string(content), "\n")
	} else if urlFlag != "" {
		urls = strings.Split(urlFlag, ",")
	} else {
		return fmt.Errorf("please provide URLs using --url or --file")
	}

	services := output.ParseServices(snoopFlag)

	ua := userAgent
	if ua == "" {
		ua = defaultUserAgent
	}

	cfg := pipeline.Config{
		Workers:      workers,
		Timeout:      timeout,
		Retries:      retries,
		Services:     services,
		Delay:        delay,
		UserAgent:    ua,
		CheckAccess:  checkAccess,
	}

	pipe, err := pipeline.New(cfg)
	if err != nil {
		return err
	}

	results := pipe.Run(urls)

	return output.Format(os.Stdout, results, outputFmt, checkAccess)
}
