# GO-BARCODE

[![Build Status](https://travis-ci.org/hyperworks/go-barcode.svg)](https://travis-ci.org/hyperworks/go-barcode)

* [zxing-cpp](https://github.com/glassechidna/zxing-cpp)
* [git-subtree](https://git-scm.com/book/en/v1/Git-Tools-Subtree-Merging)

ZXING-CPP is added using Git subtree merging.

# SETUP

GO-BARCODE requires ZXING-CPP which must be built separately on your platform of choice.
Since ZXING-CPP uses CMake, there is no simple way to integrate this with CGO so a build
script is provided instead. This script needs to be invoked once to produce a `libzxing.a`
file which CGO will then link into your application directly during `go build`.

```sh
$ go get github.com/hyperworks/go-barcode
$ cd $GOPATH/src/github.com/hyperworks/go-barcode
$ ./build-zxing.sh

...

$ # ready to use
```

**NOTE:** CMake hardcodes the current path into the generated Makefiles, so if you ever
move this repository after running the build script, you may need to clean the folder and
re-run the script again.

# CLI

To use this from the CLI, run `go install` on the CLI package:

```sh
$ go install github.com/hyperworks/go-barcode/scan
$ go install github.com/hyperworks/go-barcode/pngize
$ $GOPATH/bin/pngize barcode.pdf
barcode.pdf
  /var/folders/lb/d7m0fh9j58zby9l0nhznl7dr0000gn/T/barcode.pdf601205962/output.png
$ $GOPATH/bin/scan barcode.png
barcode.png
  9876543210128
```

# FAQ

* Error: `image: unknown format` - Adds a side-effect import for your image format such as
  `import _ "image/gif"` or `import _ "image/bmp"`.

