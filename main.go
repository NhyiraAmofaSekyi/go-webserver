package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	middleware "github.com/NhyiraAmofaSekyi/go-webserver/internal/middleware"
	v1 "github.com/NhyiraAmofaSekyi/go-webserver/internal/v1"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	start := time.Now()

	router := http.NewServeMux()

	log.Println("server running on port:", port)
	v1 := v1.NewRouter()
	api := "/api/v1/"
	router.Handle(api, http.StripPrefix(strings.TrimRight(api, "/"), v1))

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.CorsWrapper,
	)

	server := &http.Server{
		Handler: stack(router), // wrapped handler
		Addr:    ":" + port,    // Listen address
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
		// Notify that the server has stopped after ListenAndServe returns.

	}()
	healthEndpoint := "http://" + host + ":" + port + api + "healthz"
	go func() {
		for {
			resp, err := http.Get(healthEndpoint)
			if err == nil && resp.StatusCode == http.StatusOK {
				log.Println("Server is ready.")
				elapsed := time.Since(start)
				log.Printf("Server ready in %s", elapsed)
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
	fmt.Println("Shutting down server...")

	// context with a timeout for the shutdown process.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//gracefully shut down the server.
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("server shutdown failed: %v", err)
	}

	fmt.Println("Server gracefully stopped.")

}
