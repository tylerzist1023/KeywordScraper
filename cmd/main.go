package main

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
	"github.com/tylerzist1023/KeywordScraper/cmd/scraper"
)

func getLine(r *bufio.Reader) (string, error) {
	str, err := r.ReadString('\n')
	str = strings.Trim(str, "\r\n")
	str = strings.Trim(str, "\n")
	if err != nil {
		return "", err
	}
	return str, nil
}

func main() {
	fmt.Print("Type your query: ");
	
	reader := bufio.NewReader(os.Stdin)
	q, err := getLine(reader)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Type how many articles you want to scrape: ");
	article_num_str, err := getLine(reader)
	if err != nil {
		log.Fatal(err)
	}
	article_num, err := strconv.ParseInt(article_num_str, 10, 0)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("https://www.bing.com/search?q=%s", url.QueryEscape(q))
	fmt.Println(url)

	scraper.ScrapeBingForArticles(q, article_num)
}