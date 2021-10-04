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
	CollectionName:        "accounts",
	ServiceAccountKeyPath: "FirebaseCredentials.json",
	DatabaseURL:           "https://gomicroservicedatabase-default-rtdb.firebaseio.com/",
}

type FirebaseConfig struct {
	CollectionName        string
	ServiceAccountKeyPath string
	DatabaseURL           string
}

func NewFirebaseRealTimeDB(ctx context.Context, conf *FirebaseConfig) (*FirebaseRealTimeDB, error) {
	// Create firebase config
	firebaseConf := &firebase.Config{
		DatabaseURL: conf.DatabaseURL,
	}

	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile(conf.ServiceAccountKeyPath)

	// Create the app that allow us to connect to firebase Realtime PrismaDB
	app, err := firebase.NewApp(ctx, firebaseConf, opt)
	if err != nil {
		return nil, err
	}

	// Get the Realtime PrismaDB client
	db, err := app.DatabaseWithURL(ctx, firebaseConf.DatabaseURL)
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

	var errs []error
	for _, acc := range accounts {
		err := m.DB.NewRef(m.formatPath(acc.ID.String())).Set(ctx, acc)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return xerrors.Concat(errs...)
}

// Update a list of domain.Account to the MemMapStorage
func (m *FirebaseRealTimeDB) Update(ctx context.Context, accounts ...*domain.Account) error {
	// Transaction update handler: This may get invoked multiple times due to retries.
	updateAccount := func(a *domain.Account) func(tn firebaseDB.TransactionNode) (interface{}, error) {
		return func(tn firebaseDB.TransactionNode) (interface{}, error) {
			// Read the current state of the node.
			var acc domain.Account
			if err := tn.Unmarshal(&acc); err != nil {
				return nil, err
			}
			// Mutate the state in memory.
			acc = *a

			// Return the new value which will be written back to the database.
			return acc, nil
		}
	}

	var errs []error
	for _, a := range accounts {
		ref := m.DB.NewRef(m.formatPath(a.ID.String()))
		if err := ref.Transaction(ctx, updateAccount(a)); err != nil {
			errs = append(errs, err)
		}
	}
	return xerrors.Concat(errs...)
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
func (m *FirebaseRealTimeDB) SearchBy(
	ctx context.Context, searchFunc func(*domain.Account) bool,
) ([]*domain.Account, error) {
	accounts, err := m.All(ctx)
	if err != nil {
		return nil, err
	}
	var matched []*domain.Account
	for _, a := range accounts {
		if searchFunc(a) {
			matched = append(matched, a)
		}
	}
	return matched, nil
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

	var errs []error
	for _, acc := range accounts {
		err := m.DB.NewRef(fmt.Sprintf("%v/%v", m.Conf.CollectionName, acc.ID.String())).Delete(ctx)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return xerrors.Concat(errs...)
}

func (m FirebaseRealTimeDB) formatPath(child string) string {
	return fmt.Sprintf("%v/%v", m.Conf.CollectionName, child)
}
