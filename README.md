# ScrapeHTTP

ScrapeHTTP is a web scraper application built with Go. It allows users to scrape links from a specified URL up to a certain depth and search for specific strings within the links.

## Prerequisites

- Go 1.24.1 or later
- Git

## Installation

1. **Clone the repository**:

   ```sh
   git clone https://github.com/yourusername/scrapehttp.git
   cd scrapehttp
   ```

2. **Install dependencies**:

   ```sh
   go mod tidy
   ```

3. **Install `air` for live reloading**:

   ```sh
   go install github.com/cosmtrek/air@v1.27.3
   ```

   Ensure that your Go binaries directory is in your PATH. Add the following line to your shell profile file (`~/.bashrc`, `~/.zshrc`, etc.):

   ```sh
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

   Reload your shell profile:

   ```sh
   source ~/.bashrc  # or ~/.zshrc, depending on your shell
   ```

## Running the Development Server

1. **Ensure you are in the project root directory**:

   ```sh
   cd /path/to/scrapehttp
   ```

2. **Run the development server with `air`**:

   ```sh
   air
   ```

   This will start the server and watch for changes, automatically rebuilding and restarting the application.

## Project Structure

```
scrapehttp/
├── cmd/
│   └── scrapehttp/
│       └── main.go
├── internal/
│   ├── scraper/
│   │   ├── scraper.go
│   │   └── scraper_test.go
│   └── web/
│       ├── handlers.go
│       └── templates.go
├── static/
│   └── styles.css
├── templates/
│   └── index.html
├── go.mod
├── go.sum
└── .air.toml
```

- **cmd/scrapehttp/main.go**: The main entry point of the application.
- **internal/scraper/**: Contains the scraping logic.
- **internal/web/**: Contains the HTTP handlers and template rendering logic.
- **static/**: Contains static files like CSS.
- **templates/**: Contains HTML templates.
- **go.mod**: Go module file.
- **go.sum**: Go dependencies file.
- **.air.toml**: Configuration file for `air`.

## Usage

1. Open your web browser and navigate to `http://localhost:8080`.
2. Enter the URL, depth, and search string in the form.
3. Click the "Start Scrape" button to start the scraping process.
4. The results will be displayed on the page and saved to a file with a timestamped filename.

## License

This project is licensed under the MIT License.

Replace `https://github.com/yourusername/scrapehttp.git` with the actual URL of your GitHub repository. This `README.md` file provides all the necessary instructions for someone new to the project to get started, including how to install dependencies and start the development server.Replace `https://github.com/yourusername/scrapehttp.git` with the actual URL of your GitHub repository. This `README.md` file provides all the necessary instructions for someone new to the project to get started, including how to install dependencies and start the development server.