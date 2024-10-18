package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"simple-go-api/internal/api"
	"simple-go-api/internal/cache"
)

func main() {
	// Create a new request multiplexer
	mux := http.NewServeMux()

	// Init In Mem Cache
	cache := cache.NewInMemCacheMap()

	// Register the routes/handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		api.RouteHandler(w, r, cache)
	})

	// Create a server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Run the server in a goroutine
	go func() {
		fmt.Println("Starting up server on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	// Wait for a termination signal
	<-stop

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt a graceful shutdown
	fmt.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server exiting")
}
