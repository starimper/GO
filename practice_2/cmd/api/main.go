package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"practice_2/handlers"
	"practice_2/middleware"
)

func main() {
	mux := http.NewServeMux()

	handler := http.HandlerFunc(handlers.TasksHandler)

	mux.Handle("/tasks",
		middleware.LoggingMiddleware(
			middleware.AuthMiddleware(handler),
		),
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
