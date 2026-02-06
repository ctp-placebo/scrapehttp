# ScrapeHTTP

A web-based link scraper built with Go that crawls websites to find links containing specific search strings. The application provides a simple web interface to configure scraping parameters and displays results in a tabular format.

## Features

- **Recursive Link Crawling**: Scrape links up to a specified depth
- **Search Filtering**: Find links containing specific search strings
- **Web Interface**: User-friendly form-based interface with live results
- **Result Logging**: Automatically saves results to timestamped files
- **Smart Filtering**: Skips external links, mailto, tel, and javascript links
- **Result Management**: Maintains only the 5 most recent result files

## Prerequisites

- Go 1.24.1 or later

## Installation

1. **Clone the repository**:

   ```sh
   git clone https://github.com/ctp-placebo/scrapehttp.git
   cd scrapehttp
   ```

2. **Install dependencies**:

   ```sh
   go mod download
   ```

## Running the Application

### Standard Run

From the project root directory, run:

```sh
go run cmd/scrapehttp/main.go
```

The server will start on `http://localhost:8080`.

### Build and Run

To build a binary:

```sh
go build -o scrapehttp cmd/scrapehttp/main.go
./scrapehttp
```

### Development with Live Reload (Optional)

For development with automatic reloading, install and use `air`:

```sh
go install github.com/air-verse/air@latest
air
```

## Usage

1. **Access the Web Interface**: Open your browser and navigate to `http://localhost:8080`

2. **Configure Scraping Parameters**:
   - **URL**: Enter the starting URL to scrape (e.g., `https://example.com`)
   - **Depth**: Set the maximum crawl depth (0 = only the starting page, 1 = starting page + one level of links, etc.)
   - **Search String**: Enter the text to search for in link URLs

3. **Start Scraping**: Click the "Start" button to begin the scraping process

4. **View Results**: 
   - Results appear in a table showing:
     - Source URL (where the link was found)
     - Matching link URL
     - Depth at which the link was discovered
   - Results are automatically saved to `results_log/YYYY-MM-DD_HH-MM-SS_scraper-result.txt`

## Project Structure

```
scrapehttp/
├── cmd/
│   └── scrapehttp/
│       └── main.go           # Application entry point
├── internal/
│   ├── scraper/
│   │   ├── scraper.go        # Core scraping logic
│   │   └── scraper_test.go   # Tests for scraper
│   └── web/
│       ├── handlers.go       # HTTP handlers
│       └── templates.go      # Template rendering
├── static/
│   └── styles.css            # CSS styling
├── templates/
│   ├── header.html           # Page header template
│   ├── footer.html           # Page footer template
│   └── index.html            # Main page template
├── results_log/              # Generated result files (created at runtime)
├── go.mod                    # Go module definition
└── README.md
```

## How It Works

1. **Input Processing**: The web form captures the URL, depth, and search string
2. **Recursive Crawling**: The scraper visits the starting URL and recursively follows links within the same domain
3. **Filtering**: Links are filtered to:
   - Stay within the same domain
   - Match the search string
   - Exclude mailto, tel, and javascript links
4. **Result Display**: Matching links are displayed with their source URL and depth
5. **Persistence**: Results are saved to timestamped files, with automatic cleanup to keep only the 5 most recent files

## Running Tests

To run the tests for the scraper module:

```sh
go test ./internal/scraper/...
```

To run all tests with verbose output:

```sh
go test -v ./...
```

## License

This project is licensed under the MIT License.