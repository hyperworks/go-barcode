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
	json bool
}{}

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	flag.BoolVar(&config.json, "json", false, "output valid JSON instead of plain text.")
	flag.Parse()

	result := map[string][]string{}
	for _, file := range flag.Args() {
		result[file] = processFile(file)
	}

	if config.json {
		bytes, e := json.Marshal(result)
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
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	outdir, e := ioutil.TempDir("", filepath.Base(file))
	must(e)

	outfile := filepath.Join(outdir, "output.png")

	must(mw.ReadImage(file))
	must(mw.SetImageFormat("png"))
	must(mw.WriteImages(outfile, true))

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
