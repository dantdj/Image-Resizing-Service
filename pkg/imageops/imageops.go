package imageops

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
)

// Decodes the raw bytes of the image into a Golang Image object for processing.
func DecodeImage(imageBytes []byte, imageFormat string) (image.Image, error) {
	var image image.Image
	var err error

	switch imageFormat {
	case "image/png":
		image, err = png.Decode(bytes.NewReader(imageBytes))
	case "image/jpeg":
		image, err = jpeg.Decode(bytes.NewReader(imageBytes))
	default:
		err = errors.New("unsupported extension")
	}

	if err != nil {
		log.WithFields(log.Fields{
			"imageFormat": imageFormat,
			"error":       err,
		}).Error("couldn't decode content", imageFormat, err)
		return nil, err
	}

	return image, nil
}

// Encodes the provided image to a byte slice, following the format provided.
func EncodeImage(image image.Image, imageFormat string) ([]byte, error) {
	buffer := new(bytes.Buffer)

	switch imageFormat {
	case "image/png":
		encoder := &png.Encoder{
			CompressionLevel: png.DefaultCompression,
		}

		if err := encoder.Encode(buffer, image); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("failed to encode image", err)
			return nil, err
		}
	case "image/jpeg":
		options := jpeg.Options{
			Quality: jpeg.DefaultQuality,
		}
		if err := jpeg.Encode(buffer, image, &options); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("failed to encode image", err)
			return nil, err
		}
	default:
		errorString := fmt.Sprintf("couldn't find decoder for image format '%s", imageFormat)
		log.Info(errorString)
		return nil, errors.New(errorString)
	}

	return buffer.Bytes(), nil
}

// Resizes an image represented by the provided byte slice to the provided height and width.
// Will auto-detect the provided image format and return the image in the same format.
func ResizeImage(imageBytes []byte, height int, width int) ([]byte, error) {
	// `DetectContentType` detects the type of the data by examining known "magic bytes"
	imageFormat := http.DetectContentType(imageBytes)

	img, err := DecodeImage(imageBytes, imageFormat)
	if err != nil {
		return nil, err
	}

	resizedImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	encodedImage, err := EncodeImage(resizedImage, imageFormat)
	if err != nil {
		return nil, err
	}

	return encodedImage, nil
}
