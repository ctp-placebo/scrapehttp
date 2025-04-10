package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type LinkData struct {
	URL       string
	Text      string
	SourceURL string
	Depth     int
}

func CheckDeadLinks(baseURL string, maxDepth int) []LinkData {
	visited := make(map[string]bool)
	deadLinks := []LinkData{}

	var getLinks func(string, int)
	getLinks = func(currentURL string, depth int) {
		if visited[currentURL] || depth > maxDepth {
			return
		}
		visited[currentURL] = true

		resp, err := http.Get(currentURL)
		if err != nil {
			fmt.Printf("Error visiting %s: %v\n", currentURL, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Dead link found: %s\n", currentURL)
			deadLinks = append(deadLinks, LinkData{URL: currentURL, Text: currentURL, SourceURL: currentURL, Depth: depth})
			return
		}

		doc, err := html.Parse(resp.Body)
		if err != nil {
			fmt.Printf("Error parsing HTML at %s: %v\n", currentURL, err)
			return
		}

		var parseLinks func(*html.Node)
		parseLinks = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href := attr.Val

						// Skip unwanted links
						if strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") || strings.HasPrefix(href, "javascript:") {
							fmt.Printf("Skipping unwanted link: %s\n", href)
							continue
						}

						absoluteURL := resolveURL(href, baseURL)

						// Skip external links
						baseParsed, err := url.Parse(baseURL)
						linkParsed, err2 := url.Parse(absoluteURL)
						if err == nil && err2 == nil && baseParsed.Host != linkParsed.Host {
							fmt.Printf("Skipping external link: %s\n", absoluteURL)
							continue
						}

						// Check if the link is dead
						linkResp, err := http.Get(absoluteURL)
						if err != nil || linkResp.StatusCode != http.StatusOK {
							linkText := getTextContent(n)
							deadLinks = append(deadLinks, LinkData{URL: absoluteURL, Text: linkText, SourceURL: currentURL, Depth: depth})
						} else {
							getLinks(absoluteURL, depth+1)
						}
					}
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				parseLinks(c)
			}
		}
		parseLinks(doc)
	}

	getLinks(baseURL, 0)
	return deadLinks
}

// resolveURL function is defined in scraper.go

func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var textContent string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		textContent += getTextContent(c)
	}
	return strings.TrimSpace(textContent)
}
