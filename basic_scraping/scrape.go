package main

import (
	"github.com/gocolly/colly"

	"fmt"
)

const (
	LIMIT       = 10
	SEARCH_BODY = ".CpkrTDP54mqzpuCSn1Fa"
	URL_ELEM    = ".OQ_6vPwNhCeusNiEDcGp > div > div > a[href]"
	TITLE_ELEM  = ".EKtkFWMYpwzMKOYr0GYm"
	DESC_ELEM   = ".kY2IgmnCmOGjharHErah span"
)

func extract_urls(e *colly.HTMLElement) []string {
	urls := make([]string, 0)
	e.ForEach(SEARCH_BODY, func(i int, ch *colly.HTMLElement) {
		url := ch.ChildText(URL_ELEM)
		urls = append(urls, url)
	})
	return urls
}

func extract_titles(e *colly.HTMLElement) []string {
	titles := make([]string, 0)
	e.ForEach(SEARCH_BODY, func(i int, ch *colly.HTMLElement) {
		title := ch.ChildText(TITLE_ELEM)
		titles = append(titles, title)
	})
	return titles
}

func extract_descriptions(e *colly.HTMLElement) []string {
	descriptions := make([]string, 0)
	e.ForEach(SEARCH_BODY, func(i int, ch *colly.HTMLElement) {
		description := ch.ChildText(DESC_ELEM)
		descriptions = append(descriptions, description)
	})
	return descriptions
}

// It extracts url, title, and description of SERP,
// it saves results to JSON file
// it returns results in array.
func Crawl(c *colly.Collector) []SearchResult {
	results := make([]SearchResult, 0, LIMIT)
	fmt.Println("Start:")
	c.OnHTML(SEARCH_BODY, func(p *colly.HTMLElement) {

		urls := extract_urls(p)
		titles := extract_titles(p)
		descriptions := extract_descriptions(p)

		for idx := 0; idx < LIMIT; idx++ {
			results[idx].Url = urls[idx]
			results[idx].Title = titles[idx]
			results[idx].Description = descriptions[idx]
			fmt.Println(idx, ")", results[idx])
		}
	})

	return results
}
