package models

type Config struct {
	AppID              string
	FirebaseConfigJSON string

	MSClientID     string
	MSClientSecret string
	MSRedirectURI  string
	MSScopes       string

	Port string
}
