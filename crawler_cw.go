package main

import (
	"fmt"
	"net/url"
	"net/http"
	"crypto/tls"
	
	"github.com/gocolly/colly"
)

var visited_domains = make(map[string]bool)

func main() {
	fmt.Println("Give seed url and max depth:")
	var seedurl string
	var maxdepth int
	fmt.Scanln(&seedurl, &maxdepth)
	parsed_seed, err := url.Parse(seedurl)
	if err != nil {
		fmt.Println("Error parsing SEED:", err)
		return
	}
	seed_domain := parsed_seed.Host
	visited_domains[seed_domain] = true
	crawl(seedurl, maxdepth)
}

func crawl(seedurl string, maxdepth int) {
	c := colly.NewCollector(
		colly.MaxDepth(maxdepth),
	)
	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("Page Title: ", e.Text)
	})
    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        link := e.Request.AbsoluteURL(e.Attr("href"))
        parsed_link, err := url.Parse(link)
        if err != nil {
        	fmt.Println("Error parsing URL:", err)
        	return
        }
        domain := parsed_link.Host
        if link != "" && !visited_domains[domain] {
            visited_domains[domain] = true
            e.Request.Visit(link)
        }
    })
	c.OnRequest(func(r *colly.Request) {
	    fmt.Println("Crawling", r.URL)
	})
	c.WithTransport(&http.Transport{
	    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
	c.OnError(func(e *colly.Response, err error) {
	    fmt.Println("FAILED Request URL:", e.Request.URL, "\nError:", err)
	})
	err := c.Visit(seedurl)
	if err != nil {
	    fmt.Println("Error visiting page:", err)
	}
}
