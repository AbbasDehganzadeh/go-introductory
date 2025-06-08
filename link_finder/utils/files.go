package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-yaml"

	"abasd80/link_finder/models"
)

const (
	JSON  = "json"
	PLAIN = "plain"
	YAML  = "yaml"
)

type fileDataMap map[string]fileData
type fileData struct {
	models.URLResponse
	Url   string             `json:"url" yaml:"url"`
	Links []models.UrlAnchor `json:"links,omitempty" yaml:"links,omitempty"`
}

func SvaeFile(args *Arguments, urls models.UrlMap, links []models.UrlAnchor) {
	data := make(fileDataMap, 0)
	for url, values := range urls {
		d := fileData{values, url, make([]models.UrlAnchor, 0)}
		data[url] = d
	}
	for _, link := range links {
		tmpu := data[link.Url]
		tmpu.Links = append(tmpu.Links, link)
		data[link.Url] = tmpu
	}

	format := args.Format
	output := args.Output
	file, err := getFile(format, output)
	if err != nil {
		log.Panicf("File error:%v", err)
	}
	if args.Verbose {
		if file == os.Stdout {
			fmt.Println("Writhing to stadnard output...")
		} else {
			fmt.Printf("Writhing to fule (%s)...\n", file)
		}
	}
	switch format {
	case JSON:
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		err := enc.Encode(data)
		if err != nil {
			log.Panicf("Save error:%v", err)
		}
	case YAML:
		enc := yaml.NewEncoder(file)
		err := enc.Encode(data)
		if err != nil {
			log.Panicf("Save error:%v", err)
		}
	case PLAIN:
		for url, values := range data {
			if args.Verbose {
				fmt.Printf("URL >>>>\t {%+v}(%+v)\n", url, values)
			}
		}
	}
}

func getFile(format string, output string) (io.Writer, error) {
	// writable object: file, or stdout
	var file io.Writer
	filename := fmt.Sprintf("%s.%s", output, format)
	if strings.SplitN(filename, ".", 2)[0] == "" { // no save; stdout
		file = os.Stdout
		return file, nil
	}
	file, err := os.Create(filename)
	if err != nil {
		return file, err
	}
	return file, nil
}
