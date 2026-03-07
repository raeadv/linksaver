package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type WebsiteMeta struct {
	Title       string
	Description string
}

func GetWebsiteMeta(rawURL string) (*WebsiteMeta, error) {
	// Supply schema if missing
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	resp, err := http.Get(rawURL)
	if err != nil {
		// Fallback to http:// if https fails
		if strings.HasPrefix(rawURL, "https://") {
			rawURL = "http://" + strings.TrimPrefix(rawURL, "https://")
			resp, err = http.Get(rawURL)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch URL: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to fetch URL: %w", err)
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read only first 50KB — title/meta are always in <head>
	limitedBody := io.LimitReader(resp.Body, 50*1024)

	meta := &WebsiteMeta{}
	tokenizer := html.NewTokenizer(limitedBody)

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return meta, nil // EOF or error, return what we have

		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()

			switch strings.ToLower(token.Data) {
			case "title":
				if tokenizer.Next() == html.TextToken {
					meta.Title = strings.TrimSpace(tokenizer.Token().Data)
				}

			case "meta":
				name, property, content := "", "", ""
				for _, attr := range token.Attr {
					switch strings.ToLower(attr.Key) {
					case "name":
						name = strings.ToLower(attr.Val)
					case "property":
						property = strings.ToLower(attr.Val)
					case "content":
						content = attr.Val
					}
				}

				if meta.Description == "" {
					if name == "description" || property == "og:description" {
						meta.Description = strings.TrimSpace(content)
					}
				}
			}

		case html.EndTagToken:
			// Stop parsing once we're past </head>
			if tokenizer.Token().Data == "head" {
				return meta, nil
			}
		}
	}
}
