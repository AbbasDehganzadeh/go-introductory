package main

import (
	"fmt"
	"log"

	"abasd80/link_finder/db"
	"abasd80/link_finder/models"
	"abasd80/link_finder/parse"
	"abasd80/link_finder/utils"
)

type App interface {
	Crawl()
}

type AppStruct struct {
	Args      utils.Arguments
	Urlmap    *parse.UrlMap
	UrlAnchor []models.UrlAnchor
	Db        *db.DBStruct
}

func newApp() *AppStruct {
	args := utils.Arguments{
		NubmerOfWorkers: 4,
		UrlChannelGap:   100,
	}
	err := args.ParseArguments()
	if err != nil {
		log.Panicln(err)
	}
	urLMap := parse.NewUrlMap(args)
	db := db.NewDB(args)
	db.Connect()
	return &AppStruct{
		Args:      args,
		Urlmap:    urLMap,
		UrlAnchor: make([]models.UrlAnchor, 0),
		Db:        db,
	}
}

func main() {
	app := newApp()
	if app.Args.Verbose {
		fmt.Printf("App initialized ...")
	}
	app.Crawl()
	utils.SvaeFile(&app.Args, app.Urlmap.Urls, app.UrlAnchor)
	if !app.Args.Quiet {
		fmt.Println("Done!")
	}
}
