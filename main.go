package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var dev = true

type siteData struct {
	url   string
	pings int
	okays int
}

func main() {

	args := os.Args
	filename := "example.txt"

	if len(args) == 2 {
		filename = args[1]
	} else {
		fmt.Printf("No file provided. %v will be used.\n", filename)
	}

	urls := obtainUrlsFromFile(filename)
	startMonitor(urls)
}

func startMonitor(urls []string) {

	statuses := make(chan *siteData, len(urls))

	var sites []*siteData
	for _, url := range urls {
		newSite := siteData{
			url:   url,
			pings: 0,
			okays: 0,
		}
		sites = append(sites, &newSite)
	}

	for _, site := range sites {
		go checkSiteStatus(site, statuses)
	}

	for rs := range statuses {

		go func(sd *siteData) {
			if dev {
				time.Sleep(10 * time.Second)
			} else {
				time.Sleep(time.Minute)
			}

			checkSiteStatus(sd, statuses)
		}(rs)
	}
	fmt.Println("Done.")
}

func checkSiteStatus(s *siteData, c chan *siteData) {

	r, err := http.Get(s.url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.pings++
	if r.StatusCode == 200 {
		s.okays++
	}
	percentUp := 100.0 * float32(s.okays) / float32(s.pings)
	fmt.Println(r.Status, " returned from: ", s.url, " Up time percentage:", percentUp, ", ", s.okays, " OKs out of ", s.pings, "pings.")

	c <- s
}

func obtainUrlsFromFile(filename string) []string {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileContents := strings.ReplaceAll(string(data), "\r\n", "\n")
	urls := strings.Split(fileContents, "\n")
	if len(urls) == 0 {
		fmt.Printf("No urls obtained from %v", filename)
		os.Exit(1)
	}

	return urls
}
