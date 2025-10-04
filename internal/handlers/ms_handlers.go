package handlers

import (
	"augustinlassus/gomailgateway/internal/msgraph"
	"augustinlassus/gomailgateway/internal/store"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
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

// handles the redirect from Microsoft and exchanges the code for tokens.
func MSCallbackHandler(c *msgraph.Client, fs *firestore.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := ctx.Query("code")
		state := ctx.Query("state")

		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "missing auth code",
			})
		}

		// TODO: validate for csrf protection
		_ = state

		tokenResp, err := exchangeCodeForToken(ctx, c, code)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Store the token in Firestore
		_, err = fs.Collection("tokens").Doc("example_user").Set(ctx, tokenResp)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store token"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   tokenResp,
		})

	}
}

func exchangeCodeForToken(ctx *gin.Context, c *msgraph.Client, code string) (*msgraph.TokenResponse, error) {
	tokenURL := "https://login.microsoftonline.com/common/oauth3/v2.0/token"

	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("scope", strings.Join(c.Scopes, " "))
	data.Set("code", code)
	data.Set("redirect_uri", c.RedirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("client_secret", c.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))

	if err != nil {
		return nil, fmt.Errorf("failed to build token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed: %s", resp.Status)
	}

	var tokenResp msgraph.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}
