package parse

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"

	"abasd80/link_finder/models"
	"abasd80/link_finder/utils"
)

func (p *Parser) AddProccesser(id int) {
	if !p.args.Quiet {
		fmt.Printf("Start proccesser id %d,\n", id)
		defer fmt.Printf("Stop proccesser id %d,\n", id)
	}

	inbody := false
	var parse func(*html.Node, string)
	parse = func(n *html.Node, url string) {
		if n.Type == html.ElementNode && n.Data == "body" {
			inbody = !inbody
		}
		if inbody {
			if n.Type == html.ElementNode && n.Data == "a" {
				var atext, atitle, href string
				for _, attr := range n.Attr {
					if attr.Key == "title" {
						atitle = attr.Val
					}
					if attr.Key == "href" {
						href = attr.Val
						cUrl, canseen, err := extractUrl(href, p.args)
						seen := p.Urlmap.AddToQueue(cUrl)
						if seen {
							if p.args.Verbose {
								fmt.Println(cUrl, ", NoSeen;")
							}
						} else {
							p.urlSeen[cUrl] = true
							if p.args.Verbose {
								fmt.Println(cUrl, ", Seen;")
							}
						}
						_, ok := p.Urlmap.GetUrl(cUrl)
						if !ok { // not initialized
							local := IsLocal(cUrl, p.args.Url)
							webUrl := models.URLResponse{CanSeen: canseen, IsLocal: local, Error: &err}
							p.Urlmap.SetUrl(cUrl, webUrl)
						}
						atext = parseAnchor(n, atext)
						anchor := models.UrlAnchor{Page: url, Href: href, Title: atitle, Text: atext, Url: cUrl}
						*p.UrlAnchor = append(*p.UrlAnchor, anchor)
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parse(c, url)
		}
	}

	for url := range p.Urlmap.UrlTodoChan {
		s, b, e := p.Urlmap.visitUrl(url)
		if e != nil {
			p.Urlmap.SetUrlError(url, &e)
			log.Printf("Request Error: %v", e)
			if p.args.Verbose {
				fmt.Printf("Url %s, error %s\n", url, s.Reason)
			}
			p.Urlmap.Wg.Done()
		}
		if p.args.Verbose {
			fmt.Printf("Visiting %s...\n", url)
		}
		webUrl := models.URLResponse{
			Status:  s.Code,
			Reason:  s.Reason,
			Visited: true,
			CanSeen: true,
			IsLocal: IsLocal(url, p.args.Url),
			Error:   &e,
		}
		p.Urlmap.SetUrl(url, webUrl)
		if IsLocal(url, p.args.Url) {
			body := bytes.NewReader(b)
			el, err := html.Parse(body)
			if err != nil {
				p.Urlmap.Wg.Done()
				p.Urlmap.SetUrlError(url, &err)
				log.Printf("Parse Error: %v", err)
				if p.args.Verbose {
					fmt.Printf("Url %s, status %s\n", url, webUrl.Reason)
				}
			}
			parse(el, url)
			if !p.args.Quiet {
				fmt.Printf("Parsed %s...\n", url)
			}
			p.SaveSignal <- struct{}{}
			p.Urlmap.Wg.Done()
		}
	}
}

func (u *UrlMap) visitUrl(url string) (utils.Status, []byte, error) {
	time.Sleep(time.Millisecond * 250)
	status, body, err := utils.MakeRequest(url)
	if err != nil {
		return status, body, err
	}
	return status, body, err
}

func parseAnchor(n *html.Node, atext string) string {
	if n.Type == html.TextNode {
		btext := bytes.NewBufferString(atext)
		atext = string(fmt.Appendf(btext.Bytes(), "%s", n.Data))
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		atext = parseAnchor(c, atext)
	}
	return atext
}

func extractUrl(cUrl string, args utils.Arguments) (url_ string, canseen bool, err error) {
	canseen = true // default
	baseURL, _ := url.Parse(args.Url)
	URL, err := url.Parse(cUrl)
	if err != nil {
		log.Printf("URL Error: %v", err)
		url_, canseen = cUrl, false
		return
	}
	if !URL.IsAbs() { // not including hostname
		if URL.Scheme == "" { // default: https
			URL.Scheme = baseURL.Scheme
		}
		if URL.Host == "" { // default: base url host
			// ! better solution for stripping scheme
			host := baseURL.Host
			URL.Host = host
		}
	} else {
		if len(args.Hosts) != 0 {
			canseen = false
		}
		for _, host := range args.Hosts {
			if baseURL.Host == URL.Host {
				canseen = true
				break
			}
			if strings.Contains(URL.Host, host) {
				canseen = true
			}
		}
		for _, host := range args.NonHosts {
			if baseURL.Host == URL.Host || host == "**" {
				break
			}
			if strings.Contains(URL.Host, host) {
				canseen = false
			}
		}
	}
	// omit fragment && query
	if !strings.HasPrefix(URL.Path, "/") {
		URL.Path = "/" + URL.Path
	}
	url_ = fmt.Sprintf("%s://%s%s", URL.Scheme, URL.Host, URL.Path)
	if !args.Quiet {
		fmt.Printf("  %s:: %s,%s \n", cUrl, URL.Host, URL.Path)
	}
	return
}
