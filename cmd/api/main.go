package main

import (
	"augustinlassus/gomailgateway/internal/config"
	"augustinlassus/gomailgateway/internal/store"
	"context"
	"log"
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

}
