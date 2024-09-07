package main

import (
	"net/url"
)

func normalizeURL(orgurl string) (string, error) {
	parsedUrl, err := url.Parse(orgurl)
	if err != nil {
		return "", err
	}
	domain := parsedUrl.Hostname()
	path := parsedUrl.Path

	return domain + path, nil
}
