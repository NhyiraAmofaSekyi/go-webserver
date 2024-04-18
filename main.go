package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	middleware "github.com/NhyiraAmofaSekyi/go-webserver/internal/middleware"
	v1 "github.com/NhyiraAmofaSekyi/go-webserver/internal/v1"
)

func main() {

	router := http.NewServeMux()

	log.Println("server running on port 8080")
	v1 := v1.NewRouter()
	router.Handle("/v1/", http.StripPrefix("/v1", v1))

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.CorsWrapper,
	)

	server := &http.Server{
		Handler: stack(router), // wrapped handler
		Addr:    ":8080",       // Listen address
		// Other configurations like ReadTimeout, WriteTimeout, etc.
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
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
		fmt.Println("Server goroutine starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Block until a signal is received.
	<-quit
	fmt.Println("Shutting down server...")

	// Create a context with a timeout for the shutdown process.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server.
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("server shutdown failed: %v", err)
	}

	fmt.Println("Server gracefully stopped.")

}
