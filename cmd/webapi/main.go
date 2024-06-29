/*
Webapi is the executable for the main web server.
It builds a web server around APIs from `service/api`.
Webapi connects to external resources needed (database) and starts two web servers: the API web server, and the debug.
Everything is served via the API web server, except debug variables (/debug/vars) and profiler infos (pprof).

Usage:

	webapi [flags]

Flags and configurations are handled automatically by the code in `load-configuration.go`.

Return values (exit codes):

	0
		The program ended successfully (no errors, stopped by signal)

	> 0
		The program ended due to an error

Note that this program will update the schema of the database to the latest version available (embedded in the
executable during the build).
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/sirupsen/logrus"
	"wasa-photo.uniroma1.it/wasa-photo/service/api"
	"wasa-photo.uniroma1.it/wasa-photo/service/database_nosql"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// main is the program entry point. The only purpose of this function is to call run() and set the exit code if there is
// any error
func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(1)
	}
}

// run executes the program. The body of this function should perform the following steps:
// * reads the configuration
// * creates and configure the logger
// * connects to any external resources (like databases, authenticators, etc.)
// * creates an instance of the service/api package
// * starts the principal web server (using the service/api.Router.Handler() for HTTP handlers)
// * waits for any termination event: SIGTERM signal (UNIX), non-recoverable server error, etc.
// * closes the principal web server
func run() error {
	rand.Seed(time.Now().UnixNano())
	// Load Configuration and defaults
	cfg, err := loadConfiguration()
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return nil
		}
		return err
	}

	// Init logging
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	if cfg.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.Infof("application initializing")

	// Start Database
	logger.Println("initializing database support")
	// dbconn, err := sql.Open("sqlite3", cfg.DB.Filename)
	// Specify your custom credentials
    awsAccessKeyID := "ASIAWRRLGCS5U5O76AUS"
    awsSecretAccessKey := "KLVpK5o/hL/c29m1KzyD7pPoB1IZkRLpSGo2XQML"
    awsSessionToken := "IQoJb3JpZ2luX2VjEHYaCXVzLXdlc3QtMiJHMEUCIQCPNRi2/ockWu27cd7UBeEETH3gO5oHPeIV3/yykDqh/wIgRzOewHRB64kbOip0gEfvF9DrmfvGmgJVKfqFymZgomEquwIILxABGgw0NDk5ODg0MDAzMTUiDPyaU0OZ64kcs55fQiqYApvqQKKW/+vZTj28hgvOlyMMLyek8uq0AxjR3wTk5RxLv7KhDao2yZ7ek7p0Vte3jz0b5L2znr6yLz4CnuMPd9KI4A9sUPSE42EItYZvEhIzrRkmXg7caR/mtGpAa0cYd7FDJtGWTta3VVO7IeC9thH5eh3bWI5YQjLSClXm7sLCQPWvVBlaMJFg2eJ/h39f+ukY8PH6mMNVMvTt5q9Jwr3dMxZI3K9VDqQ6E18YQV+Ag2RG4rhuoPskXXrNR9oO+iDGLvY7seWH/wz2vZ2vePGFCX6VQ9s0yMTm0ZdCi6WrZu56VS0wrVyGmSLwiSwpF86QbKlomCcX0iBVctkv15NiwSYMUUh+16v0Amb1xR4l2PinGlmz5C8wvpqAtAY6nQHQoG3yP7cMfX/mIM38Ycfow5hd4BiUzvKsCwKWBKgQlqCt8QX4FKAazdNTt/8DI896spxMmHL7xuInBLvkbKZ2nyR8Sgs0We5AZGg7oVOnNv5x7YP8MmuV7VQKhMLAbw5zXTBa8hdBPttgPA9dGIRrAhnd07VwLIpogrI/cbLF8feCE32gJZBJ7cNOLSK0hlPTrR0fyTcQa76QLZNp"
    awsRegion := "us-east-1"

    // Create a custom AWS configuration with the provided credentials
    conf, err := config.LoadDefaultConfig(context.TODO(),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, awsSessionToken)),
        config.WithRegion(awsRegion),
    )
	if err != nil {
		logger.WithError(err).Error("error in creating an AWS session for dynamodb")
		return fmt.Errorf("connecting to AWS for dynamodb: %w", err)
	}
	/*
	defer func() {
		logger.Debug("database stopping")
		_ = sess.Close()
	}()
	*/
	svc := dynamodb.NewFromConfig(conf)
	
	// db, err := database.New(dbconn)
	db, err := database_nosql.New(svc)
	if err != nil {
		logger.WithError(err).Error("error creating AppDatabase")
		return fmt.Errorf("creating AppDatabase: %w", err)
	}

	// Start (main) API server
	logger.Info("initializing API server")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Create the API router
	apirouter, err := api.New(api.Config{
		Logger:   logger,
		Database: db,
	})
	if err != nil {
		logger.WithError(err).Error("error creating the API server instance")
		return fmt.Errorf("creating the API server instance: %w", err)
	}
	router := apirouter.Handler()

	router, err = registerWebUI(router)
	if err != nil {
		logger.WithError(err).Error("error registering web UI handler")
		return fmt.Errorf("registering web UI handler: %w", err)
	}

	// Apply CORS policy
	router = applyCORSHandler(router)

	// Create the API server
	apiserver := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           router,
		ReadTimeout:       cfg.Web.ReadTimeout,
		ReadHeaderTimeout: cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	// Start the service listening for requests in a separate goroutine
	go func() {
		logger.Infof("API listening on %s", apiserver.Addr)
		serverErrors <- apiserver.ListenAndServe()
		logger.Infof("stopping API server")
	}()

	// Create storage for photos
	path := "./storage"
	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			logger.WithError(err).Error("error creating storage for photos")
			return fmt.Errorf("creating storage for photos: %w", err)
		}
	}

	// Waiting for shutdown signal or POSIX signals
	select {
	case err := <-serverErrors:
		// Non-recoverable server error
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Infof("signal %v received, start shutdown", sig)

		// Asking API server to shut down and load shed.
		err := apirouter.Close()
		if err != nil {
			logger.WithError(err).Warning("graceful shutdown of apirouter error")
		}

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and load shed.
		err = apiserver.Shutdown(ctx)
		if err != nil {
			logger.WithError(err).Warning("error during graceful shutdown of HTTP server")
			err = apiserver.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
