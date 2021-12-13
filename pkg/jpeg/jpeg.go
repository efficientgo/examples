package jpeg

import (
	"image/jpeg"
	"os"

	"github.com/efficientgo/tools/core/pkg/errcapture"
)

func decodeEncode(imageURL string) (err error) {
	f, err := os.Open(imageURL)
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, f.Close, "close")
	img, err := jpeg.Decode(f)
	if err != nil {
		return err
	}

	// If we would like to edit the image, we would do it here (:

	o, err := os.Create("out.jpg")
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, o.Close, "close")

	return jpeg.Encode(o, img, &jpeg.Options{Quality: 100})
}
