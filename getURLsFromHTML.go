package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseUrl, err := url.Parse(rawBaseURL)
	if err != nil {
		return []string{}, err
	}
	r := strings.NewReader(htmlBody)
	doc, err := html.Parse(r)
	if err != nil {
		return []string{}, err
	}
	var urls []string
	var treeSearch func(*html.Node)
	treeSearch = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, att := range n.Attr {
				if att.Key == "href" {
					actualUrl, err := url.Parse(att.Val)
					if err != nil {

						fmt.Printf("couldn't parse href '%v': %v\n", att.Val, err)
						continue
					}
					resolvedUrl := baseUrl.ResolveReference(actualUrl)
					urls = append(urls, resolvedUrl.String())

					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			treeSearch(c)
		}
	}
	treeSearch(doc)
	return urls, nil
}
