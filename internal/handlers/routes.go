package handlers

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, fs *firestore.Client) {
	r.GET("healthz", HealthCheckHandler)
	r.GET("readyz", ReadyCheckHandler(fs))

	// TODO: implement handlers
	auth := r.Group("/auth")
	{
		ms := auth.Group("/ms")
		{
			ms.GET("/login")
			ms.GET("/callback")
		}
	}

	// TODO: perhaps add an api versioning mechanism
	api := r.Group("/api")
	{
		mail := api.Group("/mail")
		{
			ms := mail.Group("/ms")
			{
				// TODO: define routes
				ms.GET("/")
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
func ReadyCheckHandler(fs *firestore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// only checking firestore readiness for now...
		// cuz there aint nin else
		iter := fs.Collections(c)
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
