package main

import (
	"image"
	"log"
	"os"

	"github.com/hyperworks/go-barcode"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	for _, arg := range os.Args[1:] {
		processImage(arg)
	}
}

func processImage(path string) {
	log.Printf("processing %s", path)

	file, e := os.Open(path)
	noError(e)
	defer file.Close()

	img, _, e := image.Decode(file)
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
