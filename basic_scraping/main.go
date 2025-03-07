package main

import (
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"

	"fmt"
	"os"
)

type SearchResult struct {
	Description, Title, Url string
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Err: unable to configure url,")
	}
	URL := os.Getenv("WEB_URL")
	args := os.Args
	query := args[1]
	fmt.Println(query, args)

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Req: send request to", URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Res: get response from", URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("%s ain't ducking", URL)
	})

	Crawl(c /*Collector*/)

	q_url := fmt.Sprintf("%s?q=%v", URL, query)

	c.Visit(q_url)
	fmt.Print(q_url)
}
