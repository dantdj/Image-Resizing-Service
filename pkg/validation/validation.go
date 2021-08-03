package validation

import (
	"image"
	"io"

	log "github.com/sirupsen/logrus"
)

// Returns whether the file extension provided is one supported by the service.
// Currently requires an exact match (including the . in the extension)
func validateExtension(extension string, supportedExtensions []string) bool {
	for _, supportedExtension := range supportedExtensions {
		if extension == supportedExtension {
			return true
		}
	}
	log.WithField("extension", extension).Info("value not in supported extensions list", extension)

	return false
}

// Returns whether the dimensions of the provided image meet the limitations provided by the service
func validateImageDimensions(imageToValidate io.Reader, maxWidth int, maxHeight int) bool {
	imageConfig, _, err := image.DecodeConfig(imageToValidate)
	if err != nil {
		log.WithField("error", err).Error("error decoding image config")
		return false
	}

	if imageConfig.Width > maxWidth || imageConfig.Height > maxHeight {
		log.WithFields(log.Fields{
			"imageHeight": imageConfig.Height,
			"imageWidth":  imageConfig.Width,
		}).Info("image does not meet height or width requirements")
		return false
	}

	return true
}

// Runs validation on the image dimensions and file extension to ensure the image (and desired operation) is
// supported by the service
func ValidateImage(imageToValidate io.Reader, extension string, maxWidth int, maxHeight int, supportedExtensions []string) bool {
	return validateImageDimensions(imageToValidate, maxWidth, maxHeight) &&
		validateExtension(extension, supportedExtensions)
}
