#!/usr/bin/env bash

set -euo pipefail

./build-zxing.sh

go build -o scan-barcode ./scan
./scan-barcode barcode.png "https://raw.githubusercontent.com/hyperworks/go-barcode/master/barcode.png"

