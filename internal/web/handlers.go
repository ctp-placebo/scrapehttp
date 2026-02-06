package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"scrapehttp/internal/scraper"
)

// RenderTemplate function is declared in templates.go
type PageData struct {
	Mode          string             `json:"mode"`
	SearchResults []SearchResultData `json:"search_results"`
	DeadResults   []DeadLinkData     `json:"dead_results"`
}

type SearchResultData struct {
	PageURL   string `json:"page_url"`
	MatchText string `json:"match_text"`
	Depth     int    `json:"depth"`
}

type DeadLinkData struct {
	SourceURL string `json:"source_url"`
	BrokenURL string `json:"broken_url"`
	Depth     int    `json:"depth"`
	Status    int    `json:"status"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	baseURL := r.FormValue("url")
	maxDepth := 1
	_, err := fmt.Sscanf(r.FormValue("depth"), "%d", &maxDepth)
	if err != nil {
		http.Error(w, "Invalid depth value", http.StatusBadRequest)
		return
	}

	mode := r.FormValue("mode")
	if mode != "dead" {
		mode = "search"
	}

	searchString := r.FormValue("searchString")

	var searchResults []SearchResultData
	var deadResults []DeadLinkData

	if mode == "dead" {
		results := scraper.ScrapeDeadLinks(baseURL, maxDepth)
		deadResults = make([]DeadLinkData, 0, len(results))
		for _, result := range results {
			deadResults = append(deadResults, DeadLinkData{
				SourceURL: result.SourceURL,
				BrokenURL: result.LinkURL,
				Depth:     result.Depth,
				Status:    result.Status,
			})
		}
	} else {
		results := scraper.ScrapeSearch(baseURL, searchString, maxDepth)
		searchResults = make([]SearchResultData, 0, len(results))
		for _, result := range results {
			searchResults = append(searchResults, SearchResultData{
				PageURL:   result.PageURL,
				MatchText: result.MatchText,
				Depth:     result.Depth,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PageData{Mode: mode, SearchResults: searchResults, DeadResults: deadResults})

	// Create the results_log directory if it doesn't exist
	resultsDir := "results_log"
	if _, err := os.Stat(resultsDir); os.IsNotExist(err) {
		err := os.Mkdir(resultsDir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}
	}

	// Create a filename with the current date and time
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(resultsDir, fmt.Sprintf("%s_scraper-result.txt", currentTime))

	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	if mode == "dead" {
		for _, link := range deadResults {
			_, err := file.WriteString(fmt.Sprintf("Source: %s\nBroken: %s\nDepth: %d\nStatus: %d\n\n", link.SourceURL, link.BrokenURL, link.Depth, link.Status))
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		}
	} else {
		for _, link := range searchResults {
			_, err := file.WriteString(fmt.Sprintf("Page: %s\nMatch: %s\nDepth: %d\n\n", link.PageURL, link.MatchText, link.Depth))
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		}
	}

	// Limit the number of files to the 5 most recent
	files, err := os.ReadDir(resultsDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	// Sort files by modification time in descending order
	sort.Slice(files, func(i, j int) bool {
		infoI, _ := files[i].Info()
		infoJ, _ := files[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	// Remove older files if there are more than 5
	if len(files) > 5 {
		for _, file := range files[5:] {
			err := os.Remove(filepath.Join(resultsDir, file.Name()))
			if err != nil {
				fmt.Printf("Error removing file: %v\n", err)
				return
			}
		}
	}
}
