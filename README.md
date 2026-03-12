<p align="center">
  <h1 align="center">Snooper</h1>
</p>

<p align="center">
  <strong>A fast, reliable CLI tool to extract cloud storage links from URLs and documents</strong>
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
  <a href="#future-scope">Future Scope</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

---

## Features

- **Multiple cloud services**: Google Drive, SharePoint, Dropbox, OneDrive, Box, iCloud
- **Wide format support**: HTML, PDF, TXT, PPTX, DOCX, XLSX, ODT, ODS
- **Concurrent processing**: Configurable worker pool for parallel URL fetching
- **Reliable**: HTTP retries with exponential backoff, connection pooling, timeouts
- **Flexible output**: Text (default) or JSON for scripting

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

Snooper extracts cloud storage links from either directly provided URLs or from a file containing multiple URLs.

### Command-line Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--url` | `-u` | | Comma-separated URLs to process |
| `--file` | `-f` | | Path to a file containing URLs (one per line) |
| `--snoop` | `-s` | `drive` | Services to extract: `drive`, `sharepoint`, `dropbox`, `onedrive`, `box`, `icloud`, or `all` |
| `--workers` | `-w` | `5` | Number of concurrent workers |
| `--timeout` | `-t` | `60` | HTTP read timeout in seconds |
| `--retries` | `-r` | `3` | HTTP retry attempts |
| `--output` | `-o` | `text` | Output format: `text` or `json` |

### Examples

1. **Extract Dropbox links from a file**:
   ```bash
   ./snooper --snoop dropbox --file path/to/urls.txt
   ```

2. **Extract all cloud links from given URLs**:
   ```bash
   ./snooper --snoop all --url "https://example.com/file1.pdf","https://example.com/file2.pptx"
   ```

3. **Extract Google Drive and SharePoint links with 10 workers**:
   ```bash
   ./snooper --snoop drive,sharepoint --url "https://example.com/page.html" --workers 10
   ```

4. **JSON output for scripting**:
   ```bash
   ./snooper --snoop all --file urls.txt --output json
   ```

## Future Scope

- **Crawling support**: Extract links from nested pages (configurable depth)
- **OCR**: Optional Tesseract integration for scanned PDFs

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any changes or enhancements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
