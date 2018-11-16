package resizer

import (
	"image"
	_ "image/jpeg"
	"io"

	"github.com/disintegration/imaging"
)

type Resizer interface {
	Resize(io.Reader, int, int, io.Writer) error
}

type jpegResizer struct {
	rf imaging.ResampleFilter
}

func (r *jpegResizer) Resize(srs io.Reader, width, height int, dst io.Writer) error {
	img, _, err := image.Decode(srs)
	if err != nil {
		return err
	}

	resizedImg := imaging.Resize(img, width, height, r.rf)

	return imaging.Encode(dst, resizedImg, imaging.JPEG)
}

func NewJPEGResizer(config Config) Resizer {
	return &jpegResizer{rf: config.getFilter()}
}
