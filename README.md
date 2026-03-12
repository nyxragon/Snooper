<p align="center">
  <h1 align="center">Snooper</h1>
</p>

<p align="center">
  <strong>Fetches web pages and documents (PDF, Office, HTML), parses their content, and extracts cloud storage links embedded inside</strong>
</p>

<p align="center">
  <a href="https://github.com/nyxragon/Snooper/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go" alt="Go"></a>
  <a href="https://github.com/nyxragon/Snooper"><img src="https://img.shields.io/badge/github-nyxragon%2FSnooper-181717?logo=github" alt="GitHub"></a>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#recon-workflow">Recon Workflow</a> •
  <a href="#limitations">Limitations</a> •
  <a href="#future-scope">Future Scope</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

---

## Features

- **Multiple cloud services**: Google Drive, SharePoint, Dropbox, OneDrive, Box, iCloud
- **Parses multiple formats**: Fetches and extracts text from HTML pages, PDFs, PPTX, DOCX, XLSX, ODT, ODS, TXT
- **Concurrent processing**: Configurable worker pool for parallel fetching
- **Reliable**: HTTP retries with exponential backoff, connection pooling, timeouts
- **Flexible output**: Text (default) or JSON for scripting

## Recon Workflow

Snooper fits in the **recon phase** of bug bounty / pentesting, right after URL discovery. You feed it URLs to **pages and documents**. It fetches each one, parses the content (HTML, PDF, Office, etc.), and extracts any cloud storage links embedded inside.

```
┌─────────────────────┐     ┌─────────────────────┐     ┌─────────────────────┐
│  URL Discovery      │     │  Snooper             │     │  Manual Review      │
│  (pages, PDFs,     │────▶│  Fetch → Parse       │────▶│  Check permissions   │
│   Office docs)      │     │  → Extract links     │     │  Test for exposure   │
└─────────────────────┘     └─────────────────────┘     └─────────────────────┘
       │                              │
       │                              │
       ▼                              ▼
  waybackurls              Drive, SharePoint,
  gau, httpx               Dropbox, OneDrive,
  crawlers                 Box, iCloud links
```

### Where Snooper Sits

| Step | Tool / Phase | Output |
|------|--------------|--------|
| 1 | waybackurls, gau, crawlers | URLs to web pages, PDFs, presentations, etc. |
| 2 | **Snooper** | Fetches each URL, parses the document, extracts cloud links found in the content |
| 3 | Manual review | Misconfigs, exposed content, IDOR |

### Piping from URL Discovery Tools

Use `-` with `--file` to read URLs from stdin. Snooper will fetch each URL, parse the page or document, and extract cloud storage links from the content:

```bash
# Wayback URLs (pages, PDFs, etc.) → Snooper fetches & parses each → extracts cloud links
echo "target.com" | waybackurls | snooper --file - --snoop all

# GAU (GetAllUrls) → Snooper
echo "target.com" | gau | snooper --file - --snoop all --workers 10

# Filter live URLs first, then Snooper parses their content for cloud links
echo "target.com" | waybackurls | httpx -mc 200 -silent | snooper --file - --snoop all --output json
```

Snooper does not discover URLs itself. It takes URLs as input, fetches the content at each URL, and extracts cloud links from within that content.

## Limitations

- **No JavaScript execution** - Fetches raw HTML. Best for static content, PDFs, and Office documents. Links loaded dynamically in SPAs may not be found.
- **Unauthenticated only** - No cookies or auth headers. Targets public exposure.
- **Rate limiting** - Use `--delay` when scanning live targets to avoid WAF/ban.

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/nyxragon/Snooper.git
   ```

2. **Navigate to the project directory**:
   ```bash
   cd Snooper
   ```

3. **Build the tool**:
   ```bash
   go build -o snooper ./cmd/snooper/
   ```

## Usage

Provide URLs to web pages or documents (HTML, PDF, Office files, etc.). Snooper fetches each URL, parses the content, and extracts cloud storage links embedded inside.

### Command-line Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--url` | `-u` | | Comma-separated URLs to pages or documents to fetch and parse |
| `--file` | `-f` | | Path to a file containing URLs to pages/documents (one per line). Use `-` to read from stdin |
| `--snoop` | `-s` | `drive` | Services to extract: `drive`, `sharepoint`, `dropbox`, `onedrive`, `box`, `icloud`, or `all` |
| `--workers` | `-w` | `5` | Number of concurrent workers |
| `--timeout` | `-t` | `60` | HTTP read timeout in seconds |
| `--retries` | `-r` | `3` | HTTP retry attempts |
| `--output` | `-o` | `text` | Output format: `text` or `json` |
| `--delay` | `-d` | `0` | Delay in ms between requests per worker (use for polite scanning) |
| `--user-agent` | `-a` | Firefox-like | HTTP User-Agent header |
| `--check-access` | `-c` | `false` | HEAD each discovered link to check accessibility |

### Examples

1. **Parse documents from a URL list** (Snooper fetches each, parses content, extracts links):
   ```bash
   ./snooper --snoop dropbox --file path/to/urls.txt
   ```

2. **Parse a PDF and a presentation** (extracts cloud links from inside the documents):
   ```bash
   ./snooper --snoop all --url "https://example.com/file1.pdf","https://example.com/file2.pptx"
   ```

3. **Parse a web page** (extracts cloud links from the HTML content):
   ```bash
   ./snooper --snoop drive,sharepoint --url "https://example.com/page.html" --workers 10
   ```

4. **JSON output for scripting**:
   ```bash
   ./snooper --snoop all --file urls.txt --output json
   ```

5. **Polite scanning with delay and access check**:
   ```bash
   ./snooper --file urls.txt --snoop all --delay 200 --check-access
   ```

## Future Scope

- **Crawling support**: Extract links from nested pages (configurable depth)
- **OCR**: Optional Tesseract integration for scanned PDFs

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any changes or enhancements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
