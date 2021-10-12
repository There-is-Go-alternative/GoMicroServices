package database

import (
	"context"

	firebase "firebase.google.com/go"
	firebaseDB "firebase.google.com/go/db"
	"google.golang.org/api/option"
)

type FirebaseRealTimeDB struct {
	App  *firebase.App
	Conf *FirebaseConfig
	DB   *firebaseDB.Client
}

type FirebaseConfig struct {
	CollectionName    string
	ServiceAdsKeyPath string
	BaseConfig        *firebase.Config
}

//Initialize the database instance
func NewFirebaseRealTimeDB(ctx context.Context, conf *FirebaseConfig) (*FirebaseRealTimeDB, error) {
	opt := option.WithCredentialsFile(conf.ServiceAdsKeyPath)
	opt2 := option.WithEndpoint(conf.BaseConfig.DatabaseURL)

	app, err := firebase.NewApp(ctx, conf.BaseConfig, opt, opt2)
	if err != nil {
		return nil, err
	}
	db, err := app.DatabaseWithURL(ctx, conf.BaseConfig.DatabaseURL)
	if err != nil {
		return nil, err
	}
	return &FirebaseRealTimeDB{
		App:  app,
		Conf: conf,
		DB:   db,
	}, nil
}
