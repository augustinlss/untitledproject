package msgraph

type Client struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}
