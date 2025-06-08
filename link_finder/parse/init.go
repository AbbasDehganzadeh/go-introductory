package parse

import (
	"log"

	"abasd80/link_finder/models"
	"abasd80/link_finder/utils"
)

type Parser struct {
	args       utils.Arguments
	Urlmap     *UrlMap
	UrlAnchor  *[]models.UrlAnchor
	urlSeen    map[string]bool
	SaveSignal chan struct{}
}

func NewParser(arg utils.Arguments, um *UrlMap, ua *[]models.UrlAnchor) *Parser {
	return &Parser{
		args:       arg,
		Urlmap:     um,
		UrlAnchor:  ua,
		urlSeen:    make(map[string]bool, 0),
		SaveSignal: make(chan struct{}),
	}
}

func (p *Parser) InitParsing() {
	if p.args.Verbose {
		log.Println("Start Parsing root url!")
	}
	p.Urlmap.AddToQueue(p.args.Url)
	p.urlSeen[p.args.Url] = true
}
