package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

var (
	version string
)

type config struct {
	port                int
	env                 string
	supportedExtensions []string
	maxDimensions       struct {
		maxHeight int
		maxWidth  int
	}
}

// Define an application struct to hold the dependencies for the HTTP handlers, helpers,
// and middleware.
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	// Hard-coding these for now - in real usage it would make sense to pass this in via some JSON config
	cfg.supportedExtensions = []string{".jpeg", ".jpg", ".png"}
	cfg.maxDimensions.maxHeight = 5000
	cfg.maxDimensions.maxWidth = 5000

	log.SetFormatter(&log.JSONFormatter{})

	app := &application{
		config: cfg,
	}

	// Start the HTTP server.
	if err := app.serve(); err != nil {
		log.Fatal(err, nil)
	}
}
