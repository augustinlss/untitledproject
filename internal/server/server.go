package server

import (
	"augustinlassus/gomailgateway/internal/config"
	"net/http"

	"cloud.google.com/go/firestore"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config, storeClient *firestore.Client) (*Server, error) {
	return &Server{}, nil
}
