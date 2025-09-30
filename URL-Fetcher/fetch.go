package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

func fetch(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		return
	}
	fmt.Printf("%s", b)
}

func main() {
	var wg sync.WaitGroup

	for _, url := range os.Args[1:] {
		wg.Add(1)
		go fetch(url, &wg)
	}

	wg.Wait() // We wait until all goroutines finish
}
