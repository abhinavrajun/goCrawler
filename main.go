package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int32
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	defer cfg.mu.Unlock()
	cfg.mu.Lock()
	if val, exists := cfg.pages[normalizedURL]; exists {
		cfg.pages[normalizedURL] = val + 1
		return false
	}
	cfg.pages[normalizedURL] = 1
	return true
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()
	cfg.mu.Lock()
	lenOfPages := len(cfg.pages)
	cfg.mu.Unlock()
	if lenOfPages >= int(cfg.maxPages) {
		return

	}
	baseCurrent, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	if cfg.baseURL.Hostname() != baseCurrent.Hostname() {
		return
	}
	normalizedCurrentUrl, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	if isFirst := cfg.addPageVisit(normalizedCurrentUrl); !isFirst {
		return
	}
	htmlString, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("parsing this url : %v \n", normalizedCurrentUrl)
	urlList, err := getURLsFromHTML(htmlString, rawCurrentURL)

	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range urlList {
		cfg.wg.Add(1)
		go cfg.crawlPage(v)
	}

}

type pagesOrder struct {
	count int
	Url   string
}

func orderMaptoSlice(pages map[string]int) []pagesOrder {
	orderedPages := []pagesOrder{}
	for key, val := range pages {
		orderedPages = append(orderedPages, pagesOrder{count: val, Url: key})
	}
	for i := 0; i < len(orderedPages); i++ {
		for j := len(orderedPages) - 1; j > i; j-- {

			temp := orderedPages[j]
			if orderedPages[j].count > orderedPages[j-1].count {
				orderedPages[j] = orderedPages[j-1]
				orderedPages[j-1] = temp
			}
		}

	}
	return orderedPages
}

func printReport(pages map[string]int, baseURL string) {
	orderedPages := orderMaptoSlice(pages)

	fmt.Println("=================")
	fmt.Printf("REPORT for %v\n", baseURL)

	fmt.Println("=================")
	for _, val := range orderedPages {
		fmt.Printf("Found %v internal links to %v\n", val.count, val.Url)
	}
}

func main() {
	args := os.Args
	if len(args) < 4 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 4 {

		fmt.Println("too many arguments provided")

		os.Exit(1)

	}
	maxConcurrency, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("maxConcurrency not number")
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("maxPages not number")
		os.Exit(1)
	}
	baseDomain, err := url.Parse(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	var theStruct = config{
		pages:              map[string]int{},
		baseURL:            baseDomain,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           int32(maxPages),
	}
	fmt.Printf("starting crawl of: %v \n", args[1])
	theStruct.wg.Add(1)
	go theStruct.crawlPage(args[1])

	theStruct.wg.Wait()
	fmt.Println("final output ")
	for key, v := range theStruct.pages {
		fmt.Printf("url: %v , count: %v \n", key, v)
	}
	printReport(theStruct.pages, baseDomain.String())
	os.Exit(0)

}
