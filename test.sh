#!/usr/bin/env bash

set -euo pipefail

./build-zxing.sh

go build -o scan-barcode ./scan
./scan-barcode barcode.png

