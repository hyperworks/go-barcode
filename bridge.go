package barcode

// #cgo CPPFLAGS: -I${SRCDIR}/zxing-cpp/core/src
// #cgo darwin LDFLAGS: ${SRCDIR}/zxing-cpp/build/libzxing.a -liconv
// #cgo linux LDFLAGS: ${SRCDIR}/zxing-cpp/build/libzxing.a
// #include <stdlib.h>
// #include "bridge.h"
import "C"

import (
	"unsafe"
)

func scan(stride int, pixels []uint8) ([]string, error) {
	outputs := make([]*C.char, 8)
	count, e := C.scan(
		C.int(stride),
		C.int(len(pixels)),
		(*C.char)(unsafe.Pointer(&pixels[0])),
		C.int(len(outputs)),
		&outputs[0],
	)

	if e != nil {
		return nil, e
	}

	results := make([]string, 0, count)
	for i := 0; i < int(count); i++ {
		str := C.GoString(outputs[i])
		results = append(results, str)

		C.free(unsafe.Pointer(outputs[i]))
	}

	return results, nil
}
