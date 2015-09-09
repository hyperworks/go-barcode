package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gographics/imagick/imagick"
)

var config = struct {
	json      bool
	writesPNG bool
}{}

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	flag.BoolVar(&config.json, "json", false, "output valid JSON instead of plain text.")
	flag.BoolVar(&config.writesPNG, "png", false, "also output PNG images of each page.")
	flag.Parse()

	result := map[string][]string{}
	for _, file := range flag.Args() {
		result[file] = processFile(file)
	}

	if config.json {
		bytes, e := json.MarshalIndent(result, "", "  ")
		must(e)
		os.Stdout.Write(bytes)

	} else {
		for file, outputs := range result {
			fmt.Println(file)
			for _, output := range outputs {
				fmt.Println("\t", output)
			}
		}
	}
}

func processFile(file string) []string {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	outdir, e := ioutil.TempDir("", filepath.Base(file))
	must(e)

	// set density before reading, so PDF are read in higher density and don't cause blurry
	// images when exporting. (imagick conversion can be lossy.)
	must(mw.SetOption("density", "196"))
	must(mw.SetOption("resample", "72"))

	// split and export files
	must(mw.ReadImage(file))
	must(mw.SetCompressionQuality(100))
	must(mw.SetImageCompressionQuality(100))

	must(mw.SetImageFormat("pdf"))
	must(mw.WriteImages(filepath.Join(outdir, "page.pdf"), false))

	if config.writesPNG {
		must(mw.SetImageFormat("png"))
		must(mw.WriteImages(filepath.Join(outdir, "page.png"), false))
	}

	// accumulate results
	children, e := ioutil.ReadDir(outdir)
	must(e)

	result := make([]string, len(children))
	for i, child := range children {
		result[i] = filepath.Join(outdir, child.Name())
	}

	return result
}

func must(e error) {
	if e != nil {
		panic(e)
	}
}
