package main

import (
	"augustinlassus/gomailgateway/internal/config"
	"augustinlassus/gomailgateway/internal/server"
	"augustinlassus/gomailgateway/internal/store"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		log.Fatal("Error loading configuration: ", err)
	}

	ctx := context.Background()

	storeClient, err := store.NewClient(ctx, cfg)

	if err != nil {
		log.Fatal("Error initializing Firestore client: ", err)
	}

	defer storeClient.Close()

	srv, err := server.New(cfg, storeClient)

	if err != nil {
		log.Fatal("Error initializing server: ", err)
	}

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	log.Println("Server started on port", cfg.Port)

	// here making a channel buffer to listen to os Signal
	// made buffer size 1, meaning it can hold one Signal,
	// so channel can hold one signal before a sender must wait for a receiver.
	quit := make(chan os.Signal, 1)

	// Configure quit channel to listen for interrupt and terminate signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// this line will block until the quit channel receives a signal
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server stopped.")

}
