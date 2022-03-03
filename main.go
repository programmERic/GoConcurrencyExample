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
		fmt.Printf("No file provided. %v will be used.\n", filename)
	}

	urls := obtainUrlsFromFile(filename)
	startMonitor(urls)
}

func startMonitor(urls []string) {

	status := make(chan string)

	for i := 0; i < 2; i++ {
		for _, u := range urls {
			go checkSiteStatus(u, status)
			fmt.Println(<-status)
		}

		fmt.Println()
	}
}

func checkSiteStatus(s string, c chan string) {
	r, err := http.Get(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//time.Sleep(5 * time.Second)
	c <- fmt.Sprint(r.Status, " returned from: ", s)
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
