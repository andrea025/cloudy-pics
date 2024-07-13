// Package api exposes the main API engine. All HTTP APIs are handled here.
package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"cloudy-pics.uniroma1.it/cloudy-pics/service/database_nosql"
	"cloudy-pics.uniroma1.it/cloudy-pics/service/storage"
	"cloudy-pics.uniroma1.it/cloudy-pics/service/lambdafunc"
)

// Config is used to provide dependencies and configuration to the New function.
type Config struct {
	// Logger where log entries are sent
	Logger logrus.FieldLogger

	// Database is the instance of database_nosql.AppDatabase where data are saved
	Database database_nosql.AppDatabase

	// Storage is an instance of S3 bucket where to store images
	Storage storage.AppStorage

	// Lambda 
	Lambda lambdafunc.AppLambda
}

// Router is the package API interface representing an API handler builder
type Router interface {
	// Handler returns an HTTP handler for APIs provided in this package
	Handler() http.Handler

	// Close terminates any resource used in the package
	Close() error
}

// New returns a new Router instance
func New(cfg Config) (Router, error) {
	// Check if the configuration is correct
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Database == nil {
		return nil, errors.New("database is required")
	}
	if cfg.Storage == nil {
		return nil, errors.New("storage is required")
	}
	if cfg.Lambda == nil {
		return nil, errors.New("lambda is required")
	}

	// Create a new router where we will register HTTP endpoints. The server will pass requests to this router to be
	// handled.
	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &_router{
		router:     router,
		baseLogger: cfg.Logger,
		db:         cfg.Database,
		s3: 		cfg.Storage,
		lambda: 	cfg.Lambda,
	}, nil
}

type _router struct {
	router *httprouter.Router

	// baseLogger is a logger for non-requests contexts, like goroutines or background tasks not started by a request.
	// Use context logger if available (e.g., in requests) instead of this logger.
	baseLogger logrus.FieldLogger

	db database_nosql.AppDatabase

	s3 storage.AppStorage

	lambda lambdafunc.AppLambda
}
