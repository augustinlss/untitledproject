package msgraph

import (
	"augustinlassus/gomailgateway/internal/config"
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"golang.org/x/oauth2"
)

// Client wraps the Microsoft Graph SDK client
type Client struct {
	GraphClient     *msgraphsdk.GraphServiceClient
	Config          *config.Config
	Token           *azcore.TokenCredential
	IsDelegatedAuth bool
}

// NewClient creates a new Microsoft Graph client with client credentials flow (app-only)
func NewClient(cfg *config.Config) (*Client, error) {
	// For client credential flow, we need to use the .default scope format
	// This is required as per Microsoft's authentication requirements
	scopes := []string{"https://graph.microsoft.com/.default"}

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
		GraphClient:     graphClient,
		Config:          cfg,
		IsDelegatedAuth: false,
	}, nil
}

type oauth2TokenCredential struct {
	token *oauth2.Token
}

func (c *oauth2TokenCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     c.token.AccessToken,
		ExpiresOn: c.token.Expiry,
	}, nil
}

// NewDelegatedClient creates a Graph client authorized as a specific user
func NewDelegatedClient(cfg *config.Config, token *oauth2.Token) (*Client, error) {
	// Wrap oauth2 token in our adapter
	var cred azcore.TokenCredential = &oauth2TokenCredential{token: token}

	scopes := strings.Fields(cfg.MSScopes)
	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to create delegated Graph client: %w", err)
	}

	return &Client{
		GraphClient:     graphClient,
		Config:          cfg,
		IsDelegatedAuth: true,
	}, nil
}

// GetMe retrieves the current user's profile
// Note: This only works with delegated authentication flow, not with client credentials
func (c *Client) GetMe(ctx context.Context) (models.Userable, error) {
	if !c.IsDelegatedAuth {
		return nil, fmt.Errorf("/me request is only valid with delegated authentication flow")
	}
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
