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
	result := map[string][]string{}

	for _, item := range items {
		if u, e := url.Parse(item); e == nil && strings.HasPrefix(u.Scheme, "http") {
			result[item] = processURL(u)
		} else {
			result[item] = processPath(item)
		}
	}

	return result
}

func processURL(u *url.URL) []string {
	resp, e := http.Get(u.String())
	noError(e)

	defer resp.Body.Close()
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

	result, e := barcode.Scan(img)
	noError(e)

	if result == "" {
		return []string{}
	} else {
		return []string{result}
	}
}

func noError(e error) {
	if e != nil {
		panic(e)
	}
}
