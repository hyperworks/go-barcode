package pipeline

import (
	"io/ioutil"
	"path/filepath"
)

type File struct {
	Path  string
	Type  string
	Codes []string
}

type Pipe func(*File) []*File

func Run(input []*File, pipes ...Pipe) []*File {
	files := input
	results := []*File{}

	for _, pipe := range pipes {
		for _, file := range files {
			results = append(results, pipe(file)...)
		}

		files = results
		results = []*File{}
	}

	return files
}

func ScanAndExtract(paths ...string) map[string][]string {
	files := make([]*File, len(paths))
	for i, path := range paths {
		files[i] = &File{
			Path:  path,
			Type:  filepath.Ext(path),
			Codes: nil,
		}
	}

	files = Run(files, DownloadHTTP, ExtractPDF, Scan)
	result := map[string][]string{}
	for _, file := range files {
		result[
	}
}

func DownloadHTTP(*File) []*File {
}

func ExtractPDF(*File) []*File {
}

func Scan(*File) []*File {
}
