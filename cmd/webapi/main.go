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
    awsAccessKeyID := "ASIAWRRLGCS5VRMQC44M"
    awsSecretAccessKey := "cQ6jIrJYESBz0awMPE5MBsaqrWAS2qktrZRwMGa6"
    awsSessionToken := "IQoJb3JpZ2luX2VjEF4aCXVzLXdlc3QtMiJHMEUCIEIFYMkECDm3HwmqbHZVfRoNZKPBn4EfdM6Huk6ib7ShAiEAsqJLNGIAVSFqxvkW6h0zU33yxFfVy6d4qYr/uAkzYVkquwIIFxABGgw0NDk5ODg0MDAzMTUiDGFxWq7XsPFwBQBUjSqYAkLUx0HOx6wXEZT85Its/E1pQucMqZ7GRxhN7/KLDY0Fj+Uy47H5bhFTdM9oJllJEArsCDxRtsyy1tnvJF6pWBglcL3q1KaYIUZQRTYy9AaDPnecv3YxFYrF6vtQbM6bTn16EG5kVwPBEUjLKeknsDFExXksced7qBHz3c/enEecgK4vaD9/uNQqFiBUJP21efLcRb/YhpNX4rdxDtSWx3UZdKi671zx5MmDqxxXUVxr2fJQJYgDZTTuIxeY11qL9LdE01DReIwJ0F6B2KPdOtBamgSmVc0iyUsIO6R7L3709NzEIC9+jZ0muLI7HhhtVyBV8sfmyM/Bvmwq4FTIm+QFVEU1jGwjP5Z2WLwbMLJYCekWpl2qmGMw84H7swY6nQHhpb9IgGJ0fKAbD236ANHX4kAmonYtXQta/c8JAF+g2uar7sptyN27Thaq32Mk9/Rl2c8VRC8pqfZtZPdWAxFyK+b4vEIt0PjEFMkp+dIBL9l6MPQV28XMtyiTB+Nj9hJBrI78A8gmRM1e4MYlcm5LChZ8EZfm00XZdzKilYrHbr1skprN7dFPZqUQYj/dobWd4SrXJ9R0piyDOAtn"
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
