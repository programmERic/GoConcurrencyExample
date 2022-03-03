package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	args := os.Args
	filename := "example.txt"
	if len(args) == 2 {
		filename = args[1]
	} else {
		fmt.Println("No file provided. example.txt will be used.")
	}

	urls := obtainUrlsFromFile(filename)
	fmt.Println(urls)

	startMonitor(urls)
}

func startMonitor(urls []string) {
	for {

		for _, url := range urls {
			v := checkSiteStatus(url)

			fmt.Println(v)
		}

	}
}

func checkSiteStatus(s string) bool {
	r, err := http.Get(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return r.StatusCode != 200
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
		fmt.Errorf("No urls obtained from %v", filename)
		os.Exit(1)
	}
	return urls
}
