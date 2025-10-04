package msgraph

type Client struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}
