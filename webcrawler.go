package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/jackdanger/collectlinks"
)

//Synchronized Map
var sm sync.Map

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("No Start URL!")
		os.Exit(1)
	}
	queue := make(chan string)
	go func() {
		queue <- args[0]
	}()
	for uri := range queue {
		go retrieve(uri, queue)
	}
}

func retrieve(uri string, queue chan string) {
	sm.Store(uri, true)
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Err:", err)
		return
	}
	defer resp.Body.Close()
	links := collectlinks.All(resp.Body)
	for _, link := range links {

		absolute := fixUrl(link, uri)
		fmt.Println(absolute)
		if uri != "" {
			_, ok := sm.Load(absolute)
			if ok == false {
				go func() {
					queue <- absolute
				}()
			}
		}
	}
}

func fixUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}
