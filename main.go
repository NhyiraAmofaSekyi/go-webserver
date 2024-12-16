package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/NhyiraAmofaSekyi/go-webserver/internal/config"
	"github.com/NhyiraAmofaSekyi/go-webserver/internal/logger"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	middleware "github.com/NhyiraAmofaSekyi/go-webserver/internal/middleware"
	v1 "github.com/NhyiraAmofaSekyi/go-webserver/internal/v1"
)

func main() {

	start := time.Now()
	router := http.NewServeMux()

	config.Initialise()
	Config := config.Config
	port := strconv.Itoa(Config.APIPort)

	logger.Info("Initializing server on port: %s", port)

	v1 := v1.NewRouter()
	api := "/api/v1/"
	router.Handle(api, http.StripPrefix(strings.TrimRight(api, "/"), v1))
	router.Handle("/metrics", promhttp.Handler())

	logger.Debug("Routes configured. API path: %s", api)

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.CorsWrapper,
		middleware.Monitoring,
	)

	server := &http.Server{
		Handler: stack(router),
		Addr:    ":" + port, // Listen address
		// Other configurations like ReadTimeout, WriteTimeout, etc.
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	// Create a channel to listen for interrupt signals.
	quit := make(chan os.Signal, 1)
	// Register the given channel to receive notifications of the specified signals.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine.
	go func() {
		logger.Info("Starting server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error starting server: %v", err)
		}
	}()
	healthEndpoint := "http://" + Config.APIHost + ":" + port + api + "healthz"

	logger.Debug("Health check endpoint: %s", healthEndpoint)

	go func() {
		for {

			resp, err := http.Get(healthEndpoint)
			if err == nil && resp.StatusCode == http.StatusOK {
				elapsed := time.Since(start)
				logger.Info("Server ready in %s", elapsed)
				resp.Body.Close()
				break
			}
			if resp != nil {
				resp.Body.Close()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Block until a signal is received.
	<-quit
	logger.Info("Shutdown signal received")

	// context with a timeout for the shutdown process.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//gracefully shut down the server.
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed: %v", err)
	}

	logger.Info("Server gracefully stopped")

}
