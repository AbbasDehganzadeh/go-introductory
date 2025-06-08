package utils

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
)

type Arguments struct {
	Url      string   // base url
	Hosts    []string // list of visited hosts
	NonHosts []string // list of non-visited hosts
	Format   string   // eg: json, yaml
	Output   string
	Quiet    bool
	Verbose  bool

	NubmerOfWorkers int
	UrlChannelGap   int

	hostsStr    string
	nonHostsStr string
}

func (args *Arguments) ParseArguments() error {
	if os.Args[1][0] != '-' {
		args.Url = os.Args[1]
		os.Args = append(os.Args[:1], os.Args[2:]...)
	} else {
		flag.StringVar(&args.Url, "url", "", "bae url for crawling")
	}
	flag.StringVar(&args.hostsStr, "hosts", "", "list of hosts to be visited")
	flag.StringVar(&args.nonHostsStr, "non-hosts", "", "list of hosts to be visited")
	flag.StringVar(&args.Output, "output", "", "filename to save results")
	flag.StringVar(&args.Format, "format", "plain", "format data for result")
	flag.BoolVar(&args.Quiet, "quiet", true, "")
	flag.BoolVar(&args.Verbose, "verbose", false, "")
	flag.Parse()
	err := args.validateArguments()
	if err != nil {
		return fmt.Errorf("Arguments validation Error: %v", err)
	}
	if args.Verbose {
		fmt.Printf("Args parsed & validated: {%+v}\n", args)
		fmt.Println(os.Args)
	}
	return nil
}

func (args *Arguments) validateArguments() error {
	URL, err := url.ParseRequestURI(args.Url)
	if err != nil {
		return errors.New("base url must be valid: eg, http://example.org")
	}
	sUrl := fmt.Sprintf("%s://%s/", URL.Scheme, URL.Host) // strip path
	args.Url = sUrl
	if args.hostsStr != "" {
		args.Hosts = strings.Split(args.hostsStr, ",")
	}
	if args.nonHostsStr != "" {
		args.NonHosts = strings.Split(args.nonHostsStr, ",")
	}
	if len(args.Hosts) != 0 && len(args.NonHosts) != 0 { // either must be specified
		return errors.New("one of flags must be specified: --hosts, --non-hosts")
	}
	if args.Verbose { // either must be true
		args.Quiet = false
	}
	file := args.Output
	args.Output = strings.TrimSuffix(file, path.Ext(file))
	format := strings.ToLower(args.Format)
	if format != JSON && format != PLAIN && format != YAML { // wrong format!
		return errors.New("fomat must be correct: `json`, `yaml`")
	}
	args.Format = format
	return nil
}
