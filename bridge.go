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

func scan(stride int, pixels []uint8) (string, error) {
	result, e := C.scan(
		C.int(stride),
		C.int(len(pixels)),
		(*C.char)(unsafe.Pointer(&pixels[0])),
	)

	if e != nil {
		return "", e
	}

	str := C.GoString(result)
	C.free(unsafe.Pointer(result))
	return str, nil
}
