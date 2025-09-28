package store

import (
	"augustinlassus/gomailgateway/internal/config"
	"context"
	"encoding/base64"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func NewClient(c context.Context, cfg *config.Config) (*firestore.Client, error) {
	cfgJson, err := base64.StdEncoding.DecodeString(cfg.FirebaseConfigJSON)

	if err != nil {
		return nil, err
	}

	client, err := firestore.NewClient(c, cfg.AppID, option.WithCredentialsJSON(cfgJson))
	if err != nil {
		return nil, err
	}

	return client, nil
}
