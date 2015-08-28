#include "bridge.h"
#include <cstring>

#include <zxing/ZXing.h>
#include <zxing/common/Array.h>
#include <zxing/common/Counted.h>
#include <zxing/common/GreyscaleLuminanceSource.h>
#include <zxing/common/HybridBinarizer.h>
#include <zxing/BinaryBitmap.h>
#include <zxing/DecodeHints.h>
#include <zxing/MultiFormatReader.h>
#include <zxing/ReaderException.h>

using namespace zxing;

extern "C" {

  const char *scan(int stride, int npixels, char *pixels) {
    int height = npixels/stride;

    // creates binary bitmap from our gray pixels
    ArrayRef<char> pixelsRef(pixels, npixels);
    Ref<GreyscaleLuminanceSource> lum(new GreyscaleLuminanceSource(pixelsRef, stride, height, 0, 0, stride, height));
    Ref<HybridBinarizer> binarizer(new HybridBinarizer(lum));
    Ref<BinaryBitmap> bmp(new BinaryBitmap(binarizer));

    // actually try to scan for barcode.
    Ref<MultiFormatReader> reader(new MultiFormatReader());
    try {
      Ref<Result> result = reader->decode(bmp, DecodeHints::TRYHARDER_HINT);

      const std::string stdStr = result->getText()->getText();
      char *cstr = new char[stdStr.length()+1];
      std::strcpy(cstr, stdStr.c_str());
      return cstr;

    } catch (const ReaderException& e) {
      return NULL;

    }
  }

}
