package parse

import (
	"strings"
	"sync"

	"abasd80/link_finder/models"
	"abasd80/link_finder/utils"
)

type UrlMap struct {
	Urls        map[string]models.URLResponse
	UrlTodoChan chan string
	Mu          *sync.RWMutex
	Wg          *sync.WaitGroup
}

func NewUrlMap(args utils.Arguments) *UrlMap {
	return &UrlMap{
		Urls:        make(map[string]models.URLResponse),
		UrlTodoChan: make(chan string, args.UrlChannelGap),
		Mu:          &sync.RWMutex{},
		Wg:          &sync.WaitGroup{},
	}
}

func (u *UrlMap) GetUrl(url string) (models.URLResponse, bool) {
	u.Mu.Lock()
	res, ok := u.Urls[url]
	u.Mu.Unlock()
	return res, ok
}

func (u *UrlMap) SetUrl(url string, obj models.URLResponse) {
	u.Mu.Lock()
	obj.Saved = false
	u.Urls[url] = obj
	u.Mu.Unlock()
}

func (u *UrlMap) AddToQueue(url string) (added bool) {
	defer u.Mu.Unlock()
	u.Mu.Lock()

	webUrl, ok := u.Urls[url]
	if ok && (!webUrl.CanSeen && webUrl.Visited) {
		return
	}

	select {
	case u.UrlTodoChan <- url:
		u.Wg.Add(1)
	default:
		return
	}

	added = true
	return
}

func (u *UrlMap) SetUrlError(url string, err *error) {
	defer u.Mu.Unlock()
	u.Mu.Lock()
	obj := u.Urls[url]
	obj.Error = err
	obj.Saved = false
	u.Urls[url] = obj
}

func IsLocal(url, baseUrl string) bool {
	// url in the same origin of root
	return strings.Contains(url, baseUrl)
}
