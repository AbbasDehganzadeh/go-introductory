package main

import (
	"log"

	"abasd80/link_finder/parse"
)

func (c *AppStruct) Crawl() {
	if c.Args.Verbose {
		defer log.Println("Stop ...")
		log.Println("Start ...")
	}
	parser := parse.NewParser(c.Args, c.Urlmap, &c.UrlAnchor)

	for i := 0; i < c.Args.NubmerOfWorkers; i++ {
		go parser.AddProccesser(i + 1)
	}
	parser.InitParsing()
	go c.Db.Save(parser.SaveSignal, *c.Urlmap, &c.UrlAnchor)

	c.Urlmap.Wg.Wait() // wait for urls to complete
	close(c.Urlmap.UrlTodoChan)
	close(parser.SaveSignal)
}
