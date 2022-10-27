package main

import (
	"context"
	"go-project/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	logger := log.New(os.Stdout, "products-api ", log.LstdFlags)

	// Create the handlers
	productHandler := handlers.NewProducts(logger)

	servMux := http.NewServeMux()
	servMux.Handle("/", productHandler)

	// Create a new serve mux and register the handlers
	server := &http.Server{
		Addr:         ":9090",
		Handler:      servMux,
		ErrorLog:     logger,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		logger.Println("Staring server on port 9090")

		err := server.ListenAndServe()
		if err != nil {
			logger.Println("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	logger.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}
