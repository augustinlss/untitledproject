package handlers

import (
	"augustinlassus/gomailgateway/internal/msgraph"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

// Registers routes to main api engine.
func RegisterRoutes(r *gin.Engine, fsClient *firestore.Client, msClient *msgraph.Client) {
	r.GET("healthz", HealthCheckHandler)
	r.GET("readyz", ReadyCheckHandler(fsClient))

	// TODO: implement handlers
	auth := r.Group("/auth")
	{
		ms := auth.Group("/microsoft")
		{
			ms.GET("/login", MSLoginHandler(msClient))
			ms.GET("/callback", MSCallbackHandler(msClient, fsClient))
		}
	}

	// TODO: perhaps add an api versioning mechanism
	api := r.Group("/api")
	{
		mail := api.Group("/mail")
		{
			ms := mail.Group("/microsoft")
			{
				// Microsoft Graph mail routes
				ms.GET("/messages", GetMessagesHandler(msClient))
				ms.POST("/send", SendMailHandler(msClient))
				ms.GET("/user", GetUserInfoHandler(msClient))
			}
		}
	}

}

// Returns ok for liveness checks
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// ReadyCheckHandler
func ReadyCheckHandler(fsClient *firestore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// only checking firestore readiness for now...
		// cuz there aint nin else
		iter := fsClient.Collections(c)
		_, err := iter.Next()
		if err != nil && err.Error() != "no more items in iterator" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealty",
				"error":  err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})

	}
}
