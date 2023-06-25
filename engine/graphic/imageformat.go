package graphic

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

// ImageFormat is an image file supported format
type ImageFormat string

// List of supported image formats
const (
	JPEG ImageFormat = ".jpeg"
	JPG              = ".jpg"
	GIF              = ".gif"
	PNG              = ".png"
)

// decoders contains image decoder by supported format
var imageDecoders = map[ImageFormat]func(reader io.Reader) (image.Image, error){
	JPEG: jpeg.Decode,
	JPG:  jpeg.Decode,
	GIF:  gif.Decode,
	PNG:  png.Decode,
}
