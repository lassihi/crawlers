
package main

import (
	"fmt"
	"os"
	"net/url"
	"net/http"
	"crypto/tls"
	
	"github.com/gocolly/colly"
)

var visited_urls = make(map[string]bool)

func main() {
	var seedurl string = os.Args[1]
	if len(os.Args) == 3 {
		filename := os.Args[2]
		fmt.Println("Filename:", filename)
		f, err := os.Create(filename)
		if err != nil {
				fmt.Println(err)
				return
		}
		output := crawl(seedurl)
		l, err := f.WriteString(output)
		if err != nil {
				fmt.Println(err)
			f.Close()
				return
		}
		fmt.Println(l, "bytes written successfully!")
		err = f.Close()
		if err != nil {
				fmt.Println(err)
				return
		}
	}else {
		crawl(seedurl)
	}
}

func crawl(seedurl string) string{
	parsed_seed, err := url.Parse(seedurl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
	}
	seed_domain := parsed_seed.Host
	var result string
		
	c := colly.NewCollector(
		colly.AllowedDomains(seed_domain),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36"),
	)
	
	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("Title:", e.Text)
		result = result + "Title: " + e.Text + "\n"
	})
	
    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        link := e.Request.AbsoluteURL(e.Attr("href"))
        if link != "" && !visited_urls[link] {
            visited_urls[link] = true
            e.Request.Visit(link)
        }
    })
    
	c.OnRequest(func(r *colly.Request) {
	    fmt.Println("Crawling", r.URL)
	    result = result + "Crawling " + r.URL.String() + "\n"
	})
	
	c.WithTransport(&http.Transport{
	    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
	
	c.OnError(func(e *colly.Response, err error) {
	    fmt.Println("FAILED Request URL:", e.Request.URL, "\nError:", err)
	})
	
	err = c.Visit(seedurl)
	if err != nil {
	    fmt.Println("Error visiting page:", err)
	}

	return result
}

