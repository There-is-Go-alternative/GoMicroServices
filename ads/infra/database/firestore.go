package database

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	firebaseDB "firebase.google.com/go/db"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal/xerrors"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type FirebaseRealTimeDB struct {
	App  *firebase.App
	Conf *FirebaseConfig
	DB   *firebaseDB.Client
}

var DefaultConf = &FirebaseConfig{
	CollectionName:    "ads",
	ServiceAdsKeyPath: "FirebaseCredentials.json",
	BaseConfig: &firebase.Config{
		DatabaseURL: "https://gomicroservicedatabase-default-rtdb.firebaseio.com/",
	},
}

type FirebaseConfig struct {
	CollectionName    string
	ServiceAdsKeyPath string
	BaseConfig        *firebase.Config
}

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

// Create add list of domain.Ad to the MemMapStorage
func (m *FirebaseRealTimeDB) Create(ctx context.Context, ads ...*domain.Ad) error {
	if len(ads) == 0 {
		return nil
	}

	errs := xerrors.ErrList{}
	for _, ad := range ads {
		err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, ad.ID.String())).Set(ctx, ad)
		if err != nil {
			errs.Add(err)
		}
	}

	if !errs.Nil() {
		return errs
	}
	return nil
}

// Update a list of domain.Ad to the MemMapStorage
func (m *FirebaseRealTimeDB) Update(ctx context.Context, ads ...*domain.Ad) error {
	return nil
}

// ByID Retrieve the info that match "id".
// Strict: As ID is the key of the map, return an error if not found
func (m *FirebaseRealTimeDB) ByID(ctx context.Context, ID domain.AdID) (*domain.Ad, error) {
	var ad domain.Ad
	if err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, ID)).Get(ctx, &ad); err != nil {
		return nil, err
	}
	if ad.ID.Validate() != nil {
		return nil, errors.Wrapf(
			xerrors.ErrorWithCode{Code: xerrors.ResourceNotFound, Err: xerrors.AdNotFound}, "ID {%v}", ID,
		)
	}
	return &ad, nil
}

// SearchBy Retrieve the account that validate the search func passed as param.
func (m *FirebaseRealTimeDB) SearchBy(ctx context.Context, searchFunc func(*domain.Ad) bool) ([]*domain.Ad, error) {
	var ads []*domain.Ad
	if err := m.DB.NewRef("accounts").Child("email").Get(ctx, &ads); err != nil {
		return nil, err
	}
	return nil, nil
}

// All return all domain.Ad.
func (m *FirebaseRealTimeDB) All(ctx context.Context) ([]*domain.Ad, error) {
	var ads map[string]*domain.Ad
	if err := m.DB.NewRef(m.Conf.CollectionName).OrderByChild("id").Get(ctx, &ads); err != nil {
		return nil, err
	}
	lst := make([]*domain.Ad, 0, len(ads))
	for _, a := range ads {
		lst = append(lst, a)
	}
	return lst, nil
}

// Remove a domain.Ad from the MemMapStorage
func (m *FirebaseRealTimeDB) Remove(ctx context.Context, ads ...*domain.Ad) error {
	if len(ads) <= 0 {
		return nil
	}

	errs := xerrors.ErrList{}
	for _, ad := range ads {
		err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, ad.ID.String())).Delete(ctx)
		if err != nil {
			errs.Add(err)
		}
	}

	if !errs.Nil() {
		return errs
	}
	return nil
}
