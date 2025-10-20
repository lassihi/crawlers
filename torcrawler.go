package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"golang.org/x/net/proxy"
)

var visited_urls = make(map[string]bool)

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Println("Usage: torcrawler <url> <file>")
	} else {
		var seedurl string = os.Args[1]

		parsed_seed, err := url.Parse(seedurl)
		if err != nil {
			fmt.Println("Error parsing SEED:", err)
			return
		}
		seed_domain := parsed_seed.Host
		visited_urls[seed_domain] = true

		var outputfile *os.File
		if len(os.Args) == 3 {
			outputfile, err := os.Create(os.Args[2])
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			defer outputfile.Close()
			crawl(seedurl, outputfile)
		}
		crawl(seedurl, outputfile)
	}
}

func crawl(seedurl string, outputfile *os.File) {
	var depth int

	socks5Proxy := "127.0.0.1:9050"

	dialer, err := proxy.SOCKS5("tcp", socks5Proxy, nil, proxy.Direct)
	if err != nil {
		panic("Failed to create SOCKS5 dialer: " + err.Error())
	}

	httpTransport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
	}

	c := colly.NewCollector()

	c.WithTransport(httpTransport)

	var logline string

	c.OnRequest(func(r *colly.Request) {
		logline = "Crawling: " + r.URL.String()
		fmt.Println(logline)
	})

	c.OnError(func(e *colly.Response, err error) {
		if depth > 1 {
			fmt.Println("Error:", err)
		}
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		if outputfile != nil {
			outputfile.WriteString(logline + "\n")
		}
		logline = "Title: " + e.Text
		fmt.Println(logline)
		if outputfile != nil {
			outputfile.WriteString(logline + "\n")
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		parsed_link, err := url.Parse(link)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return
		}
		domain := parsed_link.Host
		if strings.Contains(link, ".onion") && link != "" && !visited_urls[domain] {
			visited_urls[domain] = true
			e.Request.Visit(link)
			depth++
		}
	})

	err = c.Visit(seedurl)
	if err != nil {
		fmt.Println("Error visiting page:", err)
	}
	depth++
}
