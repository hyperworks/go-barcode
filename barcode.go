package barcode

import (
	"image"
	"image/color"
	"os"
)

// Scan() scans for QR and Barcodes on the given image. Returns empty string and `nil`
// error if no code is found.
// TODO: Return ErrNoCode instead.
func Scan(img image.Image) (string, error) {
	gray := desaturate(img)
	return scan(gray.Stride, gray.Pix)
}

func desaturate(img image.Image) *image.Gray {
	b := img.Bounds()

	result := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			gray := result.ColorModel().Convert(img.At(x, y))
			result.Set(x, y, gray.(color.Gray))
		}
	}

	return result
}

// ScanFile() decodes image data from the given file path and returns the result of
// running Scan() on it.
func ScanFile(path string) (string, error) {
	file, e := os.Open(path)
	if e != nil {
		return "", e
	}

	defer file.Close()
	img, _, e := image.Decode(file)
	if e != nil {
		return "", e
	}

	return Scan(img)
}
