package db

import (
	"fmt"
	"net/url"
	"time"

	"abasd80/link_finder/models"
	"abasd80/link_finder/parse"
)

func (db *DBStruct) Save(signal chan struct{}, umap parse.UrlMap, ua *[]models.UrlAnchor) {
	if !db.args.Quiet {
		fmt.Println("~~~~ SaveDB ~~~~")
	}
	for range signal {
		umap.Mu.Lock()
		for url, values := range umap.Urls {
			if values.Saved {
				continue
			}
			if db.args.Verbose {
				fmt.Printf("> %s : (%d, %s, %v) \n", url, values.Status, values.Reason, &values.Error)
			}
			db.saveUrl(url, values)
			time.Sleep(100 * time.Millisecond)

			// mark as saved
			values.Saved = true
			umap.Urls[url] = values
		}

		for i := 0; i < len(*ua); i++ {
			if (*ua)[i].Saved {
				continue
			}
			if db.args.Verbose {
				fmt.Printf("- %s : %s\n", (*ua)[i].Page, (*ua)[i].Href)
			}
			db.saveLink((*ua)[i])
			time.Sleep(50 * time.Millisecond)

			// mark as saved
			(*ua)[i].Saved = true
		}
		umap.Mu.Unlock()
	}
}

func (db *DBStruct) saveUrl(url string, obj models.URLResponse) error {
	if obj.Status == 0 && !obj.Visited { // not visited url
		return nil
	}
	obj.Origin, obj.Path = getOriginPath(url)
	var data models.URLResponse
	err := db.repo.Get(&data, "SELECT * FROM urls WHERE origin = ? AND path = ?", obj.Origin, obj.Path)
	if err == nil {
		// update 8tem
		stmt := "UPDATE urls SET status=?, reason=?, seen=?, visit=?, error=? WHERE origin=? AND path=?"
		db.repo.MustExec(stmt, obj.Status, obj.Reason, obj.CanSeen, obj.Visited, *obj.Error, data.Origin, data.Path)
		return nil
	}
	// insert item
	stmt := "INSERT INTO urls (origin, path, status, reason, seen, visit, error) VALUES (?, ?, ?, ?, ?, ?, ?)"
	db.repo.MustExec(stmt, obj.Origin, obj.Path, obj.Status, obj.Reason, obj.CanSeen, obj.Visited, *obj.Error)
	return err
}

func (db *DBStruct) saveLink(obj models.UrlAnchor) error {
	origin, path := getOriginPath(obj.Url)
	uid, exists := db.getUrlid(origin, path)
	if !exists { // url not exists!!
		return nil
	}
	var data models.UrlAnchor
	obj.UrlId = uid
	err := db.repo.Get(&data, "SELECT * FROM links WHERE page = ? AND href = ?", obj.Page, obj.Href)
	if err == nil {
		// update 8tem
		stmt := "UPDATE links SET text=?, title=?, urlid=? where page = ? AND href = ?"
		db.repo.MustExec(stmt, obj.Text, obj.Title, obj.UrlId, data.Page, data.Href)
		return nil
	}
	// insert item
	stmt := "INSERT INTO links (page, href, text, title, urlid) VALUES (?, ?, ?, ?, ?)"
	db.repo.MustExec(stmt, obj.Page, obj.Href, obj.Text, obj.Title, obj.UrlId)
	return err
}

func (db *DBStruct) getUrlid(origin string, path string) (int, bool) {
	var id int
	err := db.repo.Get(&id, "SELECT id FROM urls WHERE origin = ? AND path = ?", origin, path)
	if err == nil { // exists
		return id, true
	}
	return 0, false
}

func getOriginPath(url_ string) (origin string, path string) {
	tmpu, _ := url.Parse(url_)
	origin = fmt.Sprintf("%s://%s", tmpu.Scheme, tmpu.Host)
	path = tmpu.Path
	return
}
