package main

import (
    "errors"
    "log"
    "os"
    "sync"
)

var fmt = log.New(os.Stdout, "", 0)

type Fetcher interface {
    // Fetch returns the body of URL and
    // a slice of URLs found on that page.
    Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
    m := map[string]bool{url: true}
    var mx sync.Mutex
    // create group to add the collection of goroutines for which we
    // need to wait
    var wg sync.WaitGroup
    var c2 func(string, int)
    c2 = func(url string, depth int) {
        // decrement the group counter after
        // execution
        defer wg.Done()
        if depth <= 0 { //we've gone as far as we should
            return
        }
        body, urls, err := fetcher.Fetch(url)
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Printf("found: %s %q\n", url, body)
        mx.Lock()
        for _, u := range urls {
            if !m[u] {
                m[u] = true
                // update the group counter
                wg.Add(1)
                go c2(u, depth-1)
            }
        }
        mx.Unlock()
    }
    wg.Add(1)
    c2(url, depth)
    wg.Wait()
}


func main() {
    Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
    body string
    urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
    if res, ok := f[url]; ok {
        return res.body, res.urls, nil
    }
    return "", nil, errors.New("not found: " + url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
    "http://golang.org/": &fakeResult{
        "The Go Programming Language",
        []string{
            "http://golang.org/pkg/",
            "http://golang.org/cmd/",
        },
    },
    "http://golang.org/pkg/": &fakeResult{
        "Packages",
        []string{
            "http://golang.org/",
            "http://golang.org/cmd/",
            "http://golang.org/pkg/fmt/",
            "http://golang.org/pkg/os/",
        },
    },
    "http://golang.org/pkg/fmt/": &fakeResult{
        "Package fmt",
        []string{
            "http://golang.org/",
            "http://golang.org/pkg/",
        },
    },
    "http://golang.org/pkg/os/": &fakeResult{
        "Package os",
        []string{
            "http://golang.org/",
            "http://golang.org/pkg/",
        },
    },
}
