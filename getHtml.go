package main

import (
	"fmt"
	"io"
	"net/http"
)

func getHTML(rawURL string) (string, error) {
	res, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode > 400 {
		return "", fmt.Errorf("Response status error")
	}
	//fmt.Println(rawURL)
	contentHeader := res.Header.Get("Content-Type")
	//fmt.Println(contentHeader)
	if contentHeader != "text/html" && contentHeader != "text/html;charset=utf-8" && contentHeader != "text/html; charset=utf-8" {
		return "", fmt.Errorf("Header not content-type text/html")
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("not able to read body, %v", err)
	}

	return string(data), nil
}
