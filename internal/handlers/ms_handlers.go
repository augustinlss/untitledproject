package handlers

import (
	"augustinlassus/gomailgateway/internal/config"
	"augustinlassus/gomailgateway/internal/msgraph"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

func MSLoginHandler(cfg *config.Config) gin.HandlerFunc {
	// prepare OAuth2 config once
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.MSClientID,
		ClientSecret: cfg.MSClientSecret,
		RedirectURL:  cfg.MSRedirectURI,
		Scopes:       strings.Split(cfg.MSScopes, " "),
		Endpoint:     microsoft.AzureADEndpoint(cfg.MSTenantID),
	}

	return func(c *gin.Context) {
		// generate state
		stateBytes := make([]byte, 16)
		_, err := rand.Read(stateBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
			return
		}
		state := base64.URLEncoding.EncodeToString(stateBytes)

		// TODO: store this state in a session or DB to validate later

		url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
		c.Redirect(http.StatusFound, url)
	}
}

func MSCallbackHandler(cfg *config.Config, fs *firestore.Client) gin.HandlerFunc {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.MSClientID,
		ClientSecret: cfg.MSClientSecret,
		RedirectURL:  cfg.MSRedirectURI,
		Scopes:       strings.Split(cfg.MSScopes, " "),
		Endpoint:     microsoft.AzureADEndpoint(cfg.MSTenantID),
	}

	return func(c *gin.Context) {
		code := c.Query("code")
		// state := c.Query("state")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
			return
		}
		// Validate `state` here (CSRF protection)

		// Exchange the code for a token
		token, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("code exchange failed: %s", err.Error())})
			return
		}

		// Create a Graph client using that token
		client, err := msgraph.NewDelegatedClient(cfg, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create delegated client"})
			return
		}

		// Call /me (can be called cus we are using delegated client)
		user, err := client.GetMe(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("GetMe failed: %s", err.Error())})
			return
		}

		// Extract fields
		uid := *user.GetId()
		mail := ""
		if user.GetMail() != nil {
			mail = *user.GetMail()
		}
		displayName := ""
		if user.GetDisplayName() != nil {
			displayName = *user.GetDisplayName()
		}

		// Store user & token in Firestore
		userData := map[string]any{
			"displayName": displayName,
			"email":       mail,
			"id":          uid,
			"loginTime":   time.Now(),
			"token":       token,
		}
		_, err = fs.Collection("users").Doc(uid).Set(context.Background(), userData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store user data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"user":    userData,
		})
	}
}
