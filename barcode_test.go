package barcode_test

import (
	"testing"

	_ "image/png"

	. "github.com/hyperworks/go-barcode"
)

const (
	TestFile = "barcode.png"
	TestCode = "9876543210128"
)

// TODO: The library works pretty well, but our test image is much too simple/easy.
func TestScan_BarcodePNG(t *testing.T) {
	results, e := ScanFile(TestFile)
	switch {
	case e != nil:
		t.Error(e)
	case len(results) != 1:
		t.Errorf("expected 1 result, got %d", len(results))
	case results[0] != TestCode:
		t.Errorf("expected %#v, got %#v", TestCode, results[0])
	}
}

func BenchmarkScan_BarcodePNG(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, e := ScanFile(TestFile)
		if e != nil {
			b.Error(e)
		}
	}
}
