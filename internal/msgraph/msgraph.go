package msgraph

import (
	"augustinlassus/gomailgateway/internal/config"
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// Client wraps the Microsoft Graph SDK client
type Client struct {
	GraphClient *msgraphsdk.GraphServiceClient
	Config      *config.Config
}

// NewClient creates a new Microsoft Graph client using the official SDK
func NewClient(cfg *config.Config) (*Client, error) {
	// Create the auth provider using client credentials
	scopes := strings.Fields(cfg.MSScopes)

	// Create credential options
	cred, err := azidentity.NewClientSecretCredential(
		cfg.MSTenantID,
		cfg.MSClientID,
		cfg.MSClientSecret,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create credentials: %v", err)
	}

	// Create the Microsoft Graph client
	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to create graph client: %v", err)
	}

	return &Client{
		GraphClient: graphClient,
		Config:      cfg,
	}, nil
}

// GetMe retrieves the current user's profile
func (c *Client) GetMe(ctx context.Context) (models.Userable, error) {
	return c.GraphClient.Me().Get(ctx, nil)
}

// GetUsers retrieves a list of users
func (c *Client) GetUsers(ctx context.Context) ([]models.Userable, error) {
	result, err := c.GraphClient.Users().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	return result.GetValue(), nil
}

// GetUserByID retrieves a specific user by ID
func (c *Client) GetUserByID(ctx context.Context, userID string) (models.Userable, error) {
	return c.GraphClient.Users().ByUserId(userID).Get(ctx, nil)
}

// SendMail sends an email from the authenticated user
func (c *Client) SendMail(ctx context.Context, subject, body string, toRecipients []string) error {
	// TODO: Implement the actual email sending logic using the Microsoft Graph API
	return fmt.Errorf("SendMail method not fully implemented - needs to be completed with proper SDK usage")
}
