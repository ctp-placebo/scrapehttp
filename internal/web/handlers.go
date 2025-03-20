package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"scrapehttp/internal/scraper"
)

type PageData struct {
	Links []LinkData `json:"links"`
}

type LinkData struct {
	URL       string `json:"url"`
	Text      string `json:"text"`
	SourceURL string `json:"source_url"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
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

	searchString := r.FormValue("searchString")

	httpLinks := scraper.ScrapeLinks(baseURL, searchString, maxDepth)

	links := make([]LinkData, 0, len(httpLinks))
	for link, source := range httpLinks {
		links = append(links, LinkData{URL: link, Text: link, SourceURL: fmt.Sprintf("%v", source)})
	}

	// Return the links as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PageData{Links: links})

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

	// Write the results to the file
	for _, link := range links {
		_, err := file.WriteString(fmt.Sprintf("Found at: %s\nLink: %s\n", link.SourceURL, link.URL))
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
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
