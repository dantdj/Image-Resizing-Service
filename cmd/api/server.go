package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func (app *application) serve() error {
	// Define the server object with some sensible timeout defaults to prevent lingering connections
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Channel to receive any errors returned by Shutdown()
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Blocks until a signal is received
		s := <-quit

		log.WithField("receivedSignal", s).Info("shutting down server")

		// Give a 5 second grace timeout for in-flight requests to finish
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		log.Info("completing background tasks...")

		// Indicate shutdown finished with no issues - we're waiting on this down below!
		shutdownError <- nil
	}()

	log.WithFields(log.Fields{
		"environment": app.config.env,
		"address":     srv.Addr,
	}).Info("starting server")

	// Calling Shutdown causes an ErrServerClosed error to be thrown - if the error
	// is anything _but_ that, then we want to return. Otherwise, proceed with shutdown
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.WithField("address", srv.Addr).Trace("stopped server")

	return nil
}
