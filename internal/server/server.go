package server

import (
	"augustinlassus/gomailgateway/internal/config"
	"augustinlassus/gomailgateway/internal/handlers"
	"augustinlassus/gomailgateway/internal/msgraph"
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
}

// Create a new server instance with the given configuration and Firestore client
func New(cfg *config.Config, storeClient *firestore.Client, msClient *msgraph.Client) (*Server, error) {
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	handlers.RegisterRoutes(router, storeClient, msClient)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}
	return &Server{httpServer: httpServer}, nil
}

// Start the server and listen for incoming requests
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Stop the server gracefully with the given context
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
