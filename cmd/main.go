package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

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

func worker(id int, terms <-chan string, h2s chan<- map[string][]string, wg *sync.WaitGroup) {
	wg.Add(1)
	for term := range terms {
		fmt.Printf("Worker %d started job for term \"%s\"\n", id, term)
		h2 := scraper.ScrapeBingForArticles(term, 5, 5)
		if len(h2) != 0 {
			h2["term"] = append(h2["term"], term)
			h2s <- h2
		}
		fmt.Printf("Worker %d finished job for term \"%s\"\n", id, term)
	}
	fmt.Printf("Worker %d is done\n", id)
	wg.Done()
}

func receiver(h2s <-chan map[string][]string, wg *sync.WaitGroup) {
	wg.Add(1)

	filename := fmt.Sprintf("headers_%d", time.Now().UnixMilli())
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	header_count, article_count := 0, 0

	file.WriteString("{ \"data\": [")
	for h2 := range h2s {
		json_raw, err := json.MarshalIndent(h2, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		json_str := strings.Replace(string(json_raw), "\\ufffd", "", -1)
		json_str = strings.Replace(json_str, "\\u0026", "&", -1)
		_, err = file.WriteString(json_str)
		for k,v := range h2 {
			if k == "term" {
				continue
			}
			article_count++
			header_count += len(v)
		}
		if err != nil {
			log.Fatal(err)
		}
		file.WriteString(",")
	}
	file.Seek(-1, 2)

	end := fmt.Sprintf("], \"article_count\":%d, \"header_count\":%d }", article_count, header_count)
	file.WriteString(end)

	wg.Done()
}

func main() {
	if len(os.Args) == 1 {
		os.Exit(0)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	term_chan := make(chan string)
	h2s_chan := make(chan map[string][]string)

	scanner := bufio.NewScanner(file)

	var wg_workers sync.WaitGroup
	var wg_receivers sync.WaitGroup
	for i := 0; i < 8; i++ {
		go worker(i, term_chan, h2s_chan, &wg_workers)
	}
	go receiver(h2s_chan, &wg_receivers)

	num_lines := 0
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSuffix(text, "\r\n")
		text = strings.TrimSuffix(text, "\n")

		if len(text) == 0 {
			continue
		}

		//fmt.Println(num_lines, text)
		term_chan <- text

		num_lines++
		fmt.Printf("Read line %d: \"%s\"\n", num_lines, text)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	close(term_chan)
	wg_workers.Wait()

	close(h2s_chan)
	wg_receivers.Wait()
}