package database

import (
	"context"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/xerrors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// MemMapStorage is a simple implementation in memory of an Account repository
//
// It uses a map to store domain.Account
type MemMapStorage struct {
	// storage keep the domain.Account in storage
	// As a map is not concurrent safe by design,
	// we use channels to communicate instructions
	storage map[string]domain.Account

	// Rcv is a channel that wait instruction to be performed.
	// Rcv is waiting instructions to perform in a go-routine launched with Init method.
	Rcv chan func() error

	// Err is a channel that send back result of requested operation.
	// Always return smth so wait for it !
	Err chan error

	// closing is a channel that stop the main loop.
	// Is Always nil until MemMapStorage.Stop is called.
	closing chan chan struct{}

	// logger is the service that manage Infos and Warning trace.
	// TODO: Change logger
	logger zerolog.Logger
	//logger *logrus.Logger
}

// NewAccountMemMapStorage return a MemMapStorage pointer initialised.
func NewAccountMemMapStorage() (m *MemMapStorage) {
	m = new(MemMapStorage)
	m.storage = make(map[string]domain.Account)
	m.Rcv = make(chan func() error)
	m.Err = make(chan error)
	m.logger = log.With().Str("service", "AccountMemMapStorage").Logger()
	m.closing = make(chan chan struct{})
	return
}

// Run Launch the underlying go routine that handle requests.
func (m *MemMapStorage) Run(ctx context.Context) error {
	m.logger.Info().Msg("Running ...")
	defer func() {
		m.cleanup()
		m.logger.Info().Msg("Run Stopped.")
	}()

	for {
		select {
		case c := <-m.Rcv:
			m.Err <- c()
		case <-ctx.Done():
			m.logger.Info().Msg("Context canceled. Stopping ...")
			return nil
		}
	}
}

// Save add list of domain.Account to the MemMapStorage, erasing existing accounts if any
// Use lambda to access the map of domain.Account
func (m *MemMapStorage) Save(accounts ...*domain.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	var errs []error
	m.Rcv <- func() error {
		for _, a := range accounts {
			if err := a.Validate(); err != nil {
				m.logger.WithLevel(zerolog.DebugLevel).Err(err).Msgf("could not add account: %v", a)
				errs = append(errs, err)
				continue
			}
			m.storage[a.ID.String()] = *a
			m.logger.Info().Interface("account", a).Msgf("Adding Account: %v", a.Email)
		}
		return xerrors.Concat(errs...)
	}
	return <-m.Err
}

// Create add list of domain.Account to the MemMapStorage
func (m *MemMapStorage) Create(_ context.Context, accounts ...*domain.Account) error {
	return m.Save(accounts...)
}

// Update a list of domain.Account to the MemMapStorage
func (m *MemMapStorage) Update(_ context.Context, accounts ...*domain.Account) error {
	return m.Save(accounts...)
}

// ByID Retrieve the info that match "id" in map of domain.Account.
// Strict: As ID is the key of the map, return an error if not found
func (m *MemMapStorage) ByID(_ context.Context, ID domain.AccountID) (*domain.Account, error) {
	// Validate AccountID requested
	if err := ID.Validate(); err == nil {
		return nil, xerrors.InvalidAccountID
	}

	// Create receive channel to hold the result.
	rcv := make(chan *domain.Account)

	// Sending the function to execute in MemMapStorage.Run
	m.Rcv <- func() error {
		//Check in map for presence of an account by its ID
		if account, ok := m.storage[ID.String()]; ok {
			rcv <- &account
			return nil
		}
		// Account not found, returning error
		rcv <- nil
		return xerrors.AccountNotFound
	}
	// Waiting for response in channels
	return <-rcv, <-m.Err
}

// SearchBy Retrieve the account that validate the search func passed as param.
func (m *MemMapStorage) SearchBy(searchFunc func(*domain.Account) bool) ([]*domain.Account, error) {
	// Create receive channel to hold the result.
	rcv := make(chan []*domain.Account)

	// Sending the function to execute in MemMapStorage.Run
	m.Rcv <- func() error {
		var accounts []*domain.Account
		for _, account := range m.storage {
			if searchFunc(&account) {
				accounts = append(accounts, &account)
			}
		}
		rcv <- accounts
		return nil
	}
	// Waiting for response in channels
	return <-rcv, <-m.Err
}

// ByEmail Retrieve the info that match "Email" in map of domain.Account.
func (m *MemMapStorage) ByEmail(_ context.Context, email string) ([]*domain.Account, error) {
	return m.SearchBy(func(a *domain.Account) bool {
		return a.Email == email
	})
}

// ByFirstname Retrieve the info that match "FirstName" in map of domain.Account.
func (m *MemMapStorage) ByFirstname(_ context.Context, firstname string) ([]*domain.Account, error) {
	return m.SearchBy(func(a *domain.Account) bool {
		return a.Firstname == firstname
	})
}

// ByLastname Retrieve the info that match "Lastname" in map of domain.Account.
func (m *MemMapStorage) ByLastname(_ context.Context, lastname string) ([]*domain.Account, error) {
	return m.SearchBy(func(a *domain.Account) bool {
		return a.Lastname == lastname
	})
}

// ByFullname Retrieve the info that match "Firstname" and "Lastname" in map of domain.Account.
func (m *MemMapStorage) ByFullname(_ context.Context, firstname, lastname string) ([]*domain.Account, error) {
	return m.SearchBy(func(a *domain.Account) bool {
		return a.Firstname == firstname && a.Lastname == lastname
	})
}

// All return all domain.Account in MemMapStorage.
// Use lambda to and a dedicated channel to access the map of domain.Account
func (m *MemMapStorage) All(_ context.Context) ([]*domain.Account, error) {
	rcv := make(chan []*domain.Account)
	m.Rcv <- func() error {
		var lst []*domain.Account
		for _, c := range m.storage {
			lst = append(lst, &c)
		}
		rcv <- lst
		return nil
	}
	return <-rcv, <-m.Err
}

// Remove a domain.Account from the MemMapStorage
// Use lambda to access the map of domain.Account
func (m *MemMapStorage) Remove(_ context.Context, accounts ...*domain.Account) error {
	if len(accounts) <= 0 {
		return nil
	}

	// Send lambda to remove commonNetwork.Client from ClientMemMapStorage
	m.Rcv <- func() error {
		for _, a := range accounts {
			delete(m.storage, a.ID.String())
		}
		return nil
	}
	return <-m.Err
}

func (m *MemMapStorage) cleanup() {
	close(m.Rcv)
	close(m.Err)
	m.storage = nil
	m.logger.WithLevel(zerolog.DebugLevel).Msg("Cleanup Done.")
}
