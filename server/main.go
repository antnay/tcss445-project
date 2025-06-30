package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/router"
	"syscall"
	"time"
)

func main() {
	r := router.NewRouter()
	defer r.PGClose()

	port := os.Getenv("SERVER_PORT")
	if len(port) == 0 {
		port = "4000"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r.Router,
	}

	go func() {
		log.Printf("Server starting on port %s\n", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// graceful stopping
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	r.PGClose()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
