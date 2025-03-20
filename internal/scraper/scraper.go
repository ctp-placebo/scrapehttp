package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ScrapeLinks(baseURL, searchString string, maxDepth int) map[string]string {
	visited := make(map[string]bool)     // To track visited URLs
	httpLinks := make(map[string]string) // To store links that match the search string

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
			fmt.Printf("Error visiting %s: HTTP %d\n", currentURL, resp.StatusCode)
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
						absoluteURL := resolveURL(href, baseURL)

						if strings.Contains(href, searchString) {
							httpLinks[href] = currentURL
						}

						baseParsed, err := url.Parse(baseURL)
						linkParsed, err2 := url.Parse(absoluteURL)
						if err == nil && err2 == nil && baseParsed.Host == linkParsed.Host {
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

	getLinks(baseURL, 0) // 2nd argument is the depth of the current URL
	return httpLinks
}

func resolveURL(href, baseURL string) string {
	parsedBase, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	parsedHref, err := url.Parse(href)
	if err != nil {
		return ""
	}
	return parsedBase.ResolveReference(parsedHref).String()
}
