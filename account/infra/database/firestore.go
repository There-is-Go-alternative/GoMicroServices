package database

import (
	"context"
	firebase "firebase.google.com/go"
	firebaseDB "firebase.google.com/go/db"
	"fmt"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/xerrors"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type FirebaseRealTimeDB struct {
	App  *firebase.App
	Conf *FirebaseConfig
	DB   *firebaseDB.Client
}

var DefaultConf = &FirebaseConfig{
	CollectionName:        "",
	ServiceAccountKeyPath: "gomicroservicedatabase-e70d6e1fd7b1.json",
	BaseConfig: &firebase.Config{
		//DatabaseURL: "https://gomicroservice-f5b5b-default-rtdb.europe-west1.firebasedatabase.app",
		DatabaseURL: "https://gomicroservicedatabase-default-rtdb.firebaseio.com/",
		//AuthOverride: &nilMap,
	},
}

type FirebaseConfig struct {
	CollectionName        string
	ServiceAccountKeyPath string
	BaseConfig            *firebase.Config
}

//type Config struct {
//	AuthOverride     *map[string]interface{} `json:"databaseAuthVariableOverride"`
//	DatabaseURL      string                  `json:"databaseURL"`
//	ProjectID        string                  `json:"projectId"`
//	ServiceAccountID string                  `json:"serviceAccountId"`
//	StorageBucket    string                  `json:"storageBucket"`
//}

func NewFirebaseRealTimeDB(ctx context.Context, conf *FirebaseConfig) (*FirebaseRealTimeDB, error) {
	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile(conf.ServiceAccountKeyPath)
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

// Create add list of domain.Account to the MemMapStorage
func (m *FirebaseRealTimeDB) Create(ctx context.Context, accounts ...*domain.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	errs := xerrors.ErrList{}
	for _, acc := range accounts {
		err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, acc.ID.String())).Set(ctx, acc)
		if err != nil {
			errs.Add(err)
		}
	}

	if !errs.Nil() {
		return errs
	}
	return nil
}

// Update a list of domain.Account to the MemMapStorage
func (m *FirebaseRealTimeDB) Update(ctx context.Context, accounts ...*domain.Account) error {
	return nil
}

// ByID Retrieve the info that match "id".
// Strict: As ID is the key of the map, return an error if not found
func (m *FirebaseRealTimeDB) ByID(ctx context.Context, ID domain.AccountID) (*domain.Account, error) {
	var acc domain.Account
	if err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, ID)).Get(ctx, &acc); err != nil {
		return nil, err
	}
	if acc.ID.Validate() != nil {
		return nil, errors.Wrapf(
			xerrors.ErrorWithCode{Code: xerrors.ResourceNotFound, Err: xerrors.AccountNotFound}, "ID {%v}", ID,
		)
	}
	return &acc, nil
}

// SearchBy Retrieve the account that validate the search func passed as param.
func (m *FirebaseRealTimeDB) SearchBy(ctx context.Context, searchFunc func(*domain.Account) bool) ([]*domain.Account, error) {
	var accounts []*domain.Account
	if err := m.DB.NewRef("accounts").Child("email").Get(ctx, &accounts); err != nil {
		return nil, err
	}
	return nil, nil
}

// ByEmail Retrieve the info that match "Email".
func (m *FirebaseRealTimeDB) ByEmail(ctx context.Context, email string) ([]*domain.Account, error) {
	return m.SearchBy(ctx, func(a *domain.Account) bool {
		return a.Email == email
	})
}

// ByFirstname Retrieve the info that match "FirstName" in map of domain.Account.
func (m *FirebaseRealTimeDB) ByFirstname(ctx context.Context, firstname string) ([]*domain.Account, error) {
	return m.SearchBy(ctx, func(a *domain.Account) bool {
		return a.Firstname == firstname
	})
}

// ByLastname Retrieve the info that match "Lastname".
func (m *FirebaseRealTimeDB) ByLastname(ctx context.Context, lastname string) ([]*domain.Account, error) {
	return m.SearchBy(ctx, func(a *domain.Account) bool {
		return a.Lastname == lastname
	})
}

// ByFullname Retrieve the info that match "Firstname" and "Lastname".
func (m *FirebaseRealTimeDB) ByFullname(ctx context.Context, firstname, lastname string) ([]*domain.Account, error) {
	return m.SearchBy(ctx, func(a *domain.Account) bool {
		return a.Firstname == firstname && a.Lastname == lastname
	})
}

// All return all domain.Account.
func (m *FirebaseRealTimeDB) All(ctx context.Context) ([]*domain.Account, error) {
	var accounts map[string]*domain.Account
	if err := m.DB.NewRef(m.Conf.CollectionName).OrderByChild("id").Get(ctx, &accounts); err != nil {
		return nil, err
	}
	lst := make([]*domain.Account, 0, len(accounts))
	for _, a := range accounts {
		lst = append(lst, a)
	}
	return lst, nil
}

// Remove a domain.Account from the MemMapStorage
func (m *FirebaseRealTimeDB) Remove(ctx context.Context, accounts ...*domain.Account) error {
	if len(accounts) <= 0 {
		return nil
	}

	errs := xerrors.ErrList{}
	for _, acc := range accounts {
		err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, acc.ID.String())).Delete(ctx)
		if err != nil {
			errs.Add(err)
		}
	}

	if !errs.Nil() {
		return errs
	}
	return nil
}
