#include "bridge.h"
#include <cstring>

#include <zxing/ZXing.h>
#include <zxing/common/Array.h>
#include <zxing/common/Counted.h>
#include <zxing/common/GreyscaleLuminanceSource.h>
#include <zxing/common/HybridBinarizer.h>
#include <zxing/multi/GenericMultipleBarcodeReader.h>
#include <zxing/BinaryBitmap.h>
#include <zxing/DecodeHints.h>
#include <zxing/MultiFormatReader.h>
#include <zxing/ReaderException.h>

using namespace std;
using namespace zxing;

extern "C" {

  int scan(int stride, int npixels, char *pixels, int maxoutput, char *outputs[]) {
    int height = npixels / stride;

    // creates binary bitmap from our gray pixels
    ArrayRef<char> pixelsRef(pixels, npixels);
    Ref<GreyscaleLuminanceSource> lum(new GreyscaleLuminanceSource(pixelsRef, stride, height, 0, 0, stride, height));
    Ref<HybridBinarizer> binarizer(new HybridBinarizer(lum));
    Ref<BinaryBitmap> bmp(new BinaryBitmap(binarizer));

    // actually try to scan for barcode.
    MultiFormatReader reader;
    Ref<multi::GenericMultipleBarcodeReader> multiReader(new multi::GenericMultipleBarcodeReader(reader));
    try {
      // if our image is at a good resolution, we do not need the TRYHARDER_HINT which can
      // slowdown scanning by atleast an order of magnitude.
      // TODO: boolean flag for toggling this?
      /* vector< Ref<Result> > results = multiReader->decodeMultiple(bmp, DecodeHints::TRYHARDER_HINT); */
      vector< Ref<Result> > results = multiReader->decodeMultiple(bmp, 0);

      int idx = 0;
      for (vector< Ref<Result> >::iterator it = results.begin();
          it != results.end() && idx < maxoutput;
          ++it, ++idx) {
        Ref<Result> result = *it;
        const std::string stdStr = result->getText()->getText();
        char *cstr = new char[stdStr.length()+1];
        std::strcpy(cstr, stdStr.c_str());
        outputs[idx] = cstr;
      }

      return idx;

    } catch (const ReaderException& e) {
      return 0;

    } catch (const IllegalArgumentException e) {
      return 0;

    }
  }

}
