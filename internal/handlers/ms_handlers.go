package handlers

import (
	"augustinlassus/gomailgateway/internal/msgraph"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// MSLoginHandler redirects the user to the Microsoft login page for OAuth2.
func MSLoginHandler(c *msgraph.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authURL, err := buildMicrosoftAuthURL(c)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.Redirect(http.StatusFound, authURL)
	}
}

// buildMicrosoftAuthURL constructs the OAuth2 authorization URL for Microsoft.
func buildMicrosoftAuthURL(c *msgraph.Client) (string, error) {
	u, err := url.Parse("https://login.microsoftonline.com/common/oauth2/v2.0/authorize")

	if err != nil {
		return "", fmt.Errorf("failed to parse auth url: %w", err)
	}

	// Generate a random state for CSRF protection
	stateBytes := make([]byte, 16)
	_, err = rand.Read(stateBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	state := base64.URLEncoding.EncodeToString(stateBytes)

	q := u.Query()
	q.Set("client_id", c.Config.MSClientID)
	q.Set("response_type", "code")
	q.Set("redirect_uri", c.Config.MSRedirectURI)
	q.Set("response_mode", "query")
	q.Set("scope", c.Config.MSScopes)
	q.Set("state", state)

	u.RawQuery = q.Encode()

	return u.String(), nil
}

// handles the redirect from Microsoft and exchanges the code for tokens.
// The callback uri is defined in the azure dashboard of the app.
func MSCallbackHandler(c *msgraph.Client, fs *firestore.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := ctx.Query("code")
		state := ctx.Query("state")

		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "missing auth code",
			})
			return
		}

		// TODO: validate state for CSRF protection
		_ = state

		// With the official SDK, we don't need to manually exchange the code
		// The SDK handles authentication through the Azure Identity library

		// Instead, we can use the client to get user information
		user, err := getUserInfo(ctx, c)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Store user info in Firestore
		userData := map[string]any{
			"displayName": *user,
			"email":       *user.GetMail(),
			"id":          *user.GetId(),
			"loginTime":   time.Now(),
		}

		_, err = fs.Collection("users").Doc(*user.GetId()).Set(ctx, userData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store user data"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"user":    userData,
		})
	}
}

// getUserInfo retrieves the user's information from Microsoft Graph
func getUserInfo(ctx context.Context, c *msgraph.Client) (*models.User, error) {
	// Use the GetMe method from our updated msgraph client
	userInterface, err := c.GetMe(ctx)
	if err != nil {
		return nil, err
	}

	// Convert the interface to a concrete User type
	user, ok := userInterface.(*models.User)
	if !ok {
		return nil, fmt.Errorf("failed to convert user interface to User type")
	}

	return user, nil
}

// GetUserInfoHandler returns the current user's profile information
func GetUserInfoHandler(c *msgraph.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := c.GetMe(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user": map[string]any{
				"displayName": *user.GetDisplayName(),
				"email":       *user.GetMail(),
				"id":          *user.GetId(),
			},
		})
	}
}

// GetMessagesHandler returns the user's messages
func GetMessagesHandler(c *msgraph.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// This would typically use the Microsoft Graph SDK to get messages
		// For example: c.GraphClient.Me().Messages().Get(ctx, nil)
		// Since we haven't implemented this method in our client yet, we'll return a placeholder

		ctx.JSON(http.StatusOK, gin.H{
			"messages": []map[string]any{
				{
					"id":      "message-1",
					"subject": "Sample Message",
					"preview": "This is a sample message",
				},
			},
		})
	}
}

// SendMailHandler sends an email using the Microsoft Graph API
func SendMailHandler(c *msgraph.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var mailRequest struct {
			Subject      string   `json:"subject" binding:"required"`
			Body         string   `json:"body" binding:"required"`
			ToRecipients []string `json:"to" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&mailRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := c.SendMail(ctx, mailRequest.Subject, mailRequest.Body, mailRequest.ToRecipients)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Email sent successfully",
		})
	}
}
