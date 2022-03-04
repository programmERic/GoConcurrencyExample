package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

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

	status := make(chan string, len(urls))

	var sites []siteData
	for _, url := range urls {
		newSite := siteData{
			url:   url,
			pings: 0,
			okays: 0,
		}
		sites = append(sites, newSite)
	}

	for i := 0; i < 2; i++ {
		for _, s := range sites {
			go checkSiteStatus(&s, status)
			fmt.Println(<-status)
		}

		fmt.Println()
	}
}

func checkSiteStatus(s *siteData, c chan string) {
	r, err := http.Get(s.url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.pings++
	if r.StatusCode == 200 {
		s.okays++
	}
	perecentUp := 100.0 * float32(s.okays) / float32(s.pings)

	time.Sleep(5 * time.Second)
	c <- fmt.Sprint(r.Status, " returned from: ", s.url, " Up time percentage:", perecentUp)
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
