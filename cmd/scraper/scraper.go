package scraper

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

func ScrapeBingForArticles(q string, article_num int, timeout_seconds int) map[string][]string {
	bing := colly.NewCollector(colly.MaxDepth(2), colly.AllowURLRevisit())
	article := colly.NewCollector(colly.Async(true), colly.AllowURLRevisit())
	link_count := 0
	h2s := make(map[string][]string)
	started := false
	url := fmt.Sprintf("https://www.bing.com/search?q=%s", url.QueryEscape(q))
	//fmt.Println(url)
	timeout := time.Now()

	var mut sync.Mutex
	article.OnHTML("html", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		e.ForEach("meta", func(_ int, e2 *colly.HTMLElement) {
			if(e2.Attr("content") == "article" && article_num > 0) {
				//fmt.Printf("Found article! URL: %s\n", url)

				e.ForEach("h2", func(_ int, e3 *colly.HTMLElement) {
					//fmt.Printf("\n%s\n", e3.Text)

					text := e3.Text
					text = strings.Join(strings.Split(text, "\xa0"), "")
					text = strings.Join(strings.Split(text, "\t"), "")
					text = strings.TrimSpace(text)
					text = strings.Join(strings.Split(text, "\ufffd"), "")
					text = strings.Replace(text, "\ufffd", "", -1)
					text = strings.Replace(text, "\\ufffd", "", -1)
					if len(text) != 0 {
						mut.Lock()
						h2s[url] = append(h2s[url], text)
						mut.Unlock()
					}					
				})

				mut.Lock()
				article_num--
				mut.Unlock()
			}
		})
	})
	article.OnResponse(func(r *colly.Response) {
		//fmt.Printf("%d Received response from URL: %s\n", link_count, r.Request.URL.String())

		mut.Lock()
		link_count--
		mut.Unlock()
	})

	bing.OnHTML("h2 a[href]", func(e *colly.HTMLElement) {
		mut.Lock()
		started = true
		link_count++
		timeout = time.Now()
		mut.Unlock()

		href := e.Attr("href")
		//fmt.Printf("%d Sending GET request to link: %s\n", link_count, href)

		err := article.Visit(href)
		if err != nil {
			log.Fatal(href, err)
		}
	})
	bing.OnResponse(func(r *colly.Response) {
		//fmt.Printf("Status Code: %d\n", r.StatusCode)
	})

	err := bing.Visit(url)
	if err != nil {
		log.Fatal(url, err)
	}

	for (!started || (article_num > 0 && link_count > 0)) && int(time.Now().Sub(timeout).Seconds()) < timeout_seconds {

	}

	return h2s
}