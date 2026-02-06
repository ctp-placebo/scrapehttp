package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type SearchResult struct {
	PageURL   string
	MatchText string
	Depth     int
}

type DeadLinkResult struct {
	SourceURL string
	LinkURL   string
	Depth     int
	Status    int
}

func ScrapeSearch(baseURL, searchString string, maxDepth int) []SearchResult {
	searchString = strings.TrimSpace(searchString)
	if searchString == "" {
		return []SearchResult{}
	}

	visited := make(map[string]bool)
	results := make([]SearchResult, 0)
	client := &http.Client{Timeout: 10 * time.Second}
	baseParsed, err := url.Parse(baseURL)
	if err != nil {
		return results
	}

	var getLinks func(string, int)
	getLinks = func(currentURL string, depth int) {
		if visited[currentURL] || depth > maxDepth {
			return
		}
		visited[currentURL] = true

		resp, err := client.Get(currentURL)
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

		var parseNode func(*html.Node)
		parseNode = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				href := ""
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href = attr.Val
						break
					}
				}

				if href != "" && !isSkippableLink(href) {
					absoluteURL := resolveURL(href, currentURL)
					if absoluteURL != "" && sameHost(baseParsed, absoluteURL) {
						if strings.Contains(href, searchString) {
							anchorText := strings.TrimSpace(extractText(n))
							matchText := anchorText
							if matchText == "" {
								matchText = href
							}
							results = append(results, SearchResult{PageURL: currentURL, MatchText: truncateText(matchText), Depth: depth})
						}
						getLinks(absoluteURL, depth+1)
					}
				}
			}

			if n.Type == html.TextNode && isSearchableTextNode(n) {
				text := strings.TrimSpace(n.Data)
				if text != "" && strings.Contains(text, searchString) {
					results = append(results, SearchResult{PageURL: currentURL, MatchText: truncateText(text), Depth: depth})
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				parseNode(c)
			}
		}
		parseNode(doc)
	}

	getLinks(baseURL, 0)
	return results
}

func ScrapeDeadLinks(baseURL string, maxDepth int) []DeadLinkResult {
	visited := make(map[string]bool)
	results := make([]DeadLinkResult, 0)
	client := &http.Client{Timeout: 10 * time.Second}
	statusCache := make(map[string]int)
	baseParsed, err := url.Parse(baseURL)
	if err != nil {
		return results
	}

	var getLinks func(string, int)
	getLinks = func(currentURL string, depth int) {
		if visited[currentURL] || depth > maxDepth {
			return
		}
		visited[currentURL] = true

		resp, err := client.Get(currentURL)
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

		var parseNode func(*html.Node)
		parseNode = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				href := ""
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href = attr.Val
						break
					}
				}

				if href != "" && !isSkippableLink(href) {
					absoluteURL := resolveURL(href, currentURL)
					if absoluteURL != "" && sameHost(baseParsed, absoluteURL) {
						status := linkStatus(client, absoluteURL, statusCache)
						if status >= 400 {
							results = append(results, DeadLinkResult{SourceURL: currentURL, LinkURL: absoluteURL, Depth: depth, Status: status})
						} else {
							getLinks(absoluteURL, depth+1)
						}
					}
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				parseNode(c)
			}
		}
		parseNode(doc)
	}

	getLinks(baseURL, 0)
	return results
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

func isSkippableLink(href string) bool {
	return href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") || strings.HasPrefix(href, "javascript:")
}

func sameHost(baseParsed *url.URL, absoluteURL string) bool {
	linkParsed, err := url.Parse(absoluteURL)
	if err != nil {
		return false
	}
	return baseParsed.Host == linkParsed.Host
}

func extractText(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			text := strings.TrimSpace(node.Data)
			if text != "" {
				if b.Len() > 0 {
					b.WriteString(" ")
				}
				b.WriteString(text)
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return b.String()
}

func isSearchableTextNode(n *html.Node) bool {
	if n.Parent == nil || n.Parent.Type != html.ElementNode {
		return true
	}
	return n.Parent.Data != "script" && n.Parent.Data != "style"
}

func truncateText(text string) string {
	const maxLen = 200
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

func linkStatus(client *http.Client, link string, cache map[string]int) int {
	if status, ok := cache[link]; ok {
		return status
	}

	resp, err := client.Head(link)
	if err != nil || resp.StatusCode == http.StatusMethodNotAllowed {
		if resp != nil {
			resp.Body.Close()
		}
		resp, err = client.Get(link)
		if err != nil {
			cache[link] = http.StatusServiceUnavailable
			return cache[link]
		}
	}
	defer resp.Body.Close()
	cache[link] = resp.StatusCode
	return resp.StatusCode
}
