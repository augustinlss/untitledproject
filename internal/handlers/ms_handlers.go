package handlers

import (
	"augustinlassus/gomailgateway/internal/msgraph"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
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
	u, err := url.Parse("https://login.micosoftonline.com/common/oauth2/v2.0/authorize")

	if err != nil {
		return "", fmt.Errorf("failed to parse auth url: %w", err)
	}

	q := u.Query()
	q.Set("client_id", c.ClientID)
	q.Set("response_type", "code")
	q.Set("redirect_uri", c.RedirectURI)
	q.Set("response_mode", "query")
	q.Set("scope", strings.Join(c.Scopes, " "))
	q.Set("state", "xyz123") // TODO: generate a random state for CSRF protection

	u.RawQuery = q.Encode()

	return u.String(), nil
}
