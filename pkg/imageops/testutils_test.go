package imageops

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

// Create a new test image and return the Image and []byte form
func newTestImage(t *testing.T, width int, height int) (image.Image, []byte) {
	upperLeft := image.Point{0, 0}
	lowerRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upperLeft, lowerRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/4 && y < height/4:
				img.Set(x, y, color.Black)
			case x >= width/4 && y >= height/4:
				img.Set(x, y, color.White)
			default:
			}
		}
	}

	buffer := new(bytes.Buffer)
	_ = jpeg.Encode(buffer, img, nil)
	testImgBytes := buffer.Bytes()

	return img, testImgBytes
}
