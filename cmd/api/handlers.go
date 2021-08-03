package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/dantdj/Image-Resizing-Server/pkg/imageops"
	"github.com/dantdj/Image-Resizing-Server/pkg/validation"
	log "github.com/sirupsen/logrus"
)

func (app *application) pingHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"serverTimestamp": time.Now().String(),
			"environment":     app.config.env,
			"version":         version,
		},
	}

	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) capabilitiesHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"capabilities": map[string][]string{
			"supportedExtensions": app.config.supportedExtensions,
		},
	}

	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		log.WithField("error", err).Error("error writing JSON response for capabilities")
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) resizeHandler(w http.ResponseWriter, r *http.Request) {
	desiredHeight := app.readInt(r.URL.Query(), "height", 0)
	desiredWidth := app.readInt(r.URL.Query(), "width", 0)

	// If both values are zero, there's no reason for this service to be called, and the
	// caller likely has an issue in their code
	if desiredHeight == 0 && desiredWidth == 0 {
		log.Info("resize endpoint given zero values for both height and width")
		app.badRequestResponse(w, r, errors.New("at least one of desired height or desired width must be non-zero"))
		return
	}

	if desiredHeight > app.config.maxDimensions.maxHeight || desiredWidth > app.config.maxDimensions.maxWidth {
		log.Info("resize endpoint given desired values too big for either height and width")
		app.badRequestResponse(w, r, errors.New("desired height or width too large"))
		return
	}

	// 10 << 20 sets a maximum upload limit of ~10MiB - anything larger will be rejected
	r.ParseMultipartForm(10 << 20)

	imageFile, handler, err := r.FormFile("image")
	if err != nil {
		log.WithField("error", err).Error("failed to read image from request")
		app.serverErrorResponse(w, r, err)
		return
	}

	extension := filepath.Ext(handler.Filename)
	if valid := validation.ValidateImage(imageFile, extension, app.config.maxDimensions.maxWidth, app.config.maxDimensions.maxHeight, app.config.supportedExtensions); !valid {
		// Extend in future with a validation struct that records the reason the image was invalid
		log.Error("image failed validation")
		app.unsupportedMediaTypeResponse(w, r)
		return
	}

	imageFile.Close()

	// TODO: Implement the method described at https://github.com/golang/go/issues/12512
	// This should avoid having to reload the image.
	reloadedImage, _, err := r.FormFile("image")
	if err != nil {
		log.WithField("error", err).Error("failed to re-read image from request")
		app.serverErrorResponse(w, r, err)
		return
	}

	defer reloadedImage.Close()

	imageBytes, err := ioutil.ReadAll(reloadedImage)
	if err != nil {
		log.WithField("error", err).Error("failed to read image bytes from file")
		app.serverErrorResponse(w, r, err)
		return
	}

	resizedImageBytes, err := imageops.ResizeImage(imageBytes, desiredHeight, desiredWidth)
	if err != nil {
		log.WithField("error", err).Error("failed to resize image")
		app.serverErrorResponse(w, r, err)
	}

	w.Write(resizedImageBytes)
}
