package config

import (
	"errors"
	"log"
	"os"
)

type Config struct {
	AppID              string
	FirebaseConfigJSON string

	MSClientID     string
	MSClientSecret string
	MSRedirectURI  string
	MSScopes       string
	MSTenantID     string

	Port string
}

func Load() (*Config, error) {
	log.Println("Loading configuration from environment variables...")

	cfg := &Config{
		AppID:              os.Getenv("__app_id"),
		FirebaseConfigJSON: os.Getenv("__firebase_config"),

		MSClientID:     os.Getenv("MS_APP_ID"),
		MSClientSecret: os.Getenv("MS_APP_SECRET"),
		MSRedirectURI:  os.Getenv("MS_REDIRECT_URI"),
		MSScopes:       os.Getenv("MS_SCOPES"),
		MSTenantID:     os.Getenv("MS_TENANT_ID"),

		Port: os.Getenv("SERVER_PORT"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	// Validate required fields
	if cfg.AppID == "" {
		return nil, errors.New("missing required env: __app_id")
	}
	if cfg.FirebaseConfigJSON == "" {
		return nil, errors.New("missing required env: __firebase_config")
	}
	if cfg.MSClientID == "" {
		return nil, errors.New("missing required env: MS_APP_ID")
	}
	if cfg.MSClientSecret == "" {
		return nil, errors.New("missing required env: MS_APP_SECRET")
	}
	if cfg.MSTenantID == "" {
		return nil, errors.New("missing required env: MS_TENANT_ID")
	}

	println("Configuration loaded successfully.")

	return cfg, nil

}
