package scraper

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func ScrapeBingForArticles(q string, article_num int) map[string][]string {
	bing := colly.NewCollector(colly.MaxDepth(2))
	article := colly.NewCollector(colly.Async(true))
	link_count := 0
	h2s := make(map[string][]string)

	url := fmt.Sprintf("https://www.bing.com/search?q=%s", url.QueryEscape(q))
	fmt.Println(url)

	var mut sync.Mutex
	article.OnHTML("html", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		e.ForEach("meta", func(_ int, e2 *colly.HTMLElement) {
			if(e2.Attr("content") == "article" && article_num > 0) {
				fmt.Printf("Found article! URL: %s\n", url)

				e.ForEach("h2", func(_ int, e3 *colly.HTMLElement) {
					fmt.Printf("\n%s\n", e3.Text)
					mut.Lock()
					h2s[url] = append(h2s[url], e3.Text)
					mut.Unlock()
				})

				mut.Lock()
				article_num--
				mut.Unlock()
			}
		})
	})
	article.OnResponse(func(r *colly.Response) {
		fmt.Printf("%d Received response from URL: %s\n", link_count, r.Request.URL.String())

		mut.Lock()
		link_count--
		mut.Unlock()
	})

	started := false

	bing.OnHTML("h2 a[href]", func(e *colly.HTMLElement) {
		mut.Lock()
		started = true
		link_count++
		mut.Unlock()

		href := e.Attr("href")
		fmt.Printf("%d Sending GET request to link: %s\n", link_count, href)

		err := article.Visit(href)
		if err != nil {
			log.Fatal(err)
		}
	})
	bing.OnResponse(func(r *colly.Response) {
		fmt.Printf("Status Code: %d\n", r.StatusCode)
	})

	err = bing.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: timeout
	for !started || (article_num > 0 && link_count > 0) {

	}
}