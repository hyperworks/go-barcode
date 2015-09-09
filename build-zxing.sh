#!/usr/bin/env bash

set -euo pipefail

pushd .
cd zxing-cpp
mkdir -p build/
cd build

cmake -DCMAKE_BUILD_TYPE=RELEASE -G "Unix Makefiles" ..
make libzxing/fast
popd
