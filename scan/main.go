package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/hyperworks/go-barcode"
)

func main() {
	jsonResult := false
	flag.BoolVar(&jsonResult, "json", false, "output result as JSON instead of text.")
	flag.Parse()

	result := processArgs(flag.Args())
	if jsonResult {
		bytes, e := json.MarshalIndent(result, "", "  ")
		noError(e)
		os.Stdout.Write(bytes)

	} else {
		for key, codes := range result {
			fmt.Println(key)
			for _, code := range codes {
				fmt.Println("\t" + code)
			}
		}
	}
}

func processArgs(items []string) map[string][]string {
	type scanResult struct {
		key   string
		codes []string
	}

	results := make(chan *scanResult)
	defer close(results)

	for _, item := range items {
		go func(key string) {
			var codes []string
			if u, e := url.Parse(item); e == nil && strings.HasPrefix(u.Scheme, "http") {
				codes = processURL(u)
			} else {
				codes = processPath(key)
			}

			results <- &scanResult{key, codes}
		}(item)
	}

	result := map[string][]string{}
	for _ = range items {
		scanned := <-results
		result[scanned.key] = scanned.codes
	}

	return result
}

func processURL(u *url.URL) []string {
	resp, e := http.Get(u.String())
	noError(e)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil
	}

	return process(resp.Body)
}

func processPath(path string) []string {
	file, e := os.Open(path)
	noError(e)

	defer file.Close()
	return process(file)
}

func process(reader io.Reader) []string {
	img, _, e := image.Decode(reader)
	noError(e)

	results, e := barcode.Scan(img)
	noError(e)
	return results
}

func noError(e error) {
	if e != nil {
		panic(e)
	}
}
