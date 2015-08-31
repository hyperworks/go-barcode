package main

import (
	"image"
	"io"
	"log"
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
	for _, arg := range os.Args[1:] {
		if u, e := url.Parse(arg); e == nil && strings.HasPrefix(u.Scheme, "http") {
			processURL(u)
		} else {
			processPath(arg)
		}
	}
}

func processURL(u *url.URL) {
	log.Printf("fetch: %s", u)
	resp, e := http.Get(u.String())
	noError(e)

	defer resp.Body.Close()
	process(resp.Body)
}

func processPath(path string) {
	log.Printf("read: %s", path)

	file, e := os.Open(path)
	noError(e)
	defer file.Close()

	process(file)
}

func process(reader io.Reader) {
	img, _, e := image.Decode(reader)
	noError(e)

	result, e := barcode.Scan(img)
	noError(e)

	if result == "" {
		log.Println("no code found.")
	} else {
		log.Println("code:", result)
	}
}

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
