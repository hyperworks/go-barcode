package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

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
	/*func TempDir(dir, prefix string) (name string, err error)*/
}

func processFile(file string) []string {
	sceneRE := regexp.MustCompile(`Scene: \d+ of (\d+)`)

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// calculate output path.
	outdir, e := ioutil.TempDir("", filepath.Base(file))
	must(e)

	must(mw.ReadImage(file))

	// extract image metadata
	info := mw.IdentifyImage()
	matches := sceneRE.FindStringSubmatch(info)
	pagesCount := 0
	if len(matches) > 1 {
		if n, e := strconv.ParseInt(matches[1], 10, 32); e != nil {
			must(e)
		} else {
			pagesCount = int(n)
		}
	}

	fmt.Println("number of pages:", pagesCount)

	// save converted image files
	infiles := make([]string, pagesCount)
	outfiles := make([]string, pagesCount)
	pngfiles := make([]string, pagesCount)
	for i := range outfiles {
		pageNum := strconv.FormatInt(int64(i), 10)
		infiles[i] = file + "[" + pageNum + "]"
		outfiles[i] = filepath.Join(outdir, "page-"+pageNum+".pdf")
		pngfiles[i] = filepath.Join(outdir, "page-"+pageNum+".png")
	}

	for i, infile := range infiles {
		must(mw.ReadImage(infile))
		must(mw.SetCompressionQuality(100))
		must(mw.SetImageCompressionQuality(100))

		must(mw.SetImageFormat("pdf"))
		must(mw.WriteImage(outfiles[i]))
		must(mw.SetImageFormat("png"))
		must(mw.WriteImage(pngfiles[i]))
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
