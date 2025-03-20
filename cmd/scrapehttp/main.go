package main

import (
	"fmt"
	"net/http"
	"os"

	"scrapehttp/internal/web"
)

func main() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", web.IndexHandler)
	http.HandleFunc("/scrape", web.ScrapeHandler)

	fmt.Println("Starting server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
