package database

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal/xerrors"
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
	storage map[string]*domain.Ad

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
	// TODO: rework doc above
	logger zerolog.Logger
}

// NewClientMemMapStorage return a MemMapStorage pointer initialised.
func NewClientMemMapStorage() (m *MemMapStorage) {
	m = new(MemMapStorage)
	m.storage = make(map[string]*domain.Ad)
	m.Rcv = make(chan func() error)
	m.Err = make(chan error)
	m.logger = log.With().Str("service", "AdMemMapStorage").Logger()
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

// Save add a domain.Account to the MemMapStorage
// Use lambda to access the map of domain.Account
func (m *MemMapStorage) Save(ads ...*domain.Ad) error {
	if len(ads) == 0 {
		return nil
	}
	m.Rcv <- func() error {
		for _, a := range ads {
			if a == nil {
				m.logger.Warn().Msg("Tried to add a null ad")
				continue
			}
			m.storage[a.ID.String()] = a
			m.logger.Info().Msgf("Adding ad: %v", a)
		}
		return nil
	}
	return <-m.Err
}

// ByID Retrieve the info that match "id" in map of domain.Account.
func (m *MemMapStorage) ByID(ID domain.AdID) (*domain.Ad, error) {
	// Validate AccountID requested
	if !ID.Validate() {
		return nil, xerrors.InvalidAdID
	}

	// Setup a receive channel to hold the result.
	rcv := make(chan *domain.Ad)

	// Sending the function to execute in MemMapStorage.Run
	m.Rcv <- func() error {
		//Check in map for presence of an account by its ID
		if ad, ok := m.storage[ID.String()]; ok {
			rcv <- ad
			return nil
		}
		// Account not found, returning error
		rcv <- nil
		return xerrors.AdNotFound
	}
	// Waiting for response in channels
	return <-rcv, <-m.Err
}

// All return all commonNetwork.Client in ClientMemMapStorage.
// Use lambda to and a dedicated channel to access the map of commonNetwork.Client
func (m *MemMapStorage) All() ([]*domain.Ad, error) {
	rcv := make(chan []*domain.Ad)
	m.Rcv <- func() error {
		var lst []*domain.Ad
		for _, c := range m.storage {
			lst = append(lst, c)
		}
		rcv <- lst
		return nil
	}
	return <-rcv, <-m.Err
}

// Remove remove a commonNetwork.Client from the ClientMemMapStorage
// Use lambda to access the map of commonNetwork.Client
func (m *MemMapStorage) Remove(ads ...*domain.Ad) error {
	if len(ads) <= 0 {
		return nil
	}

	// Send lambda to remove commonNetwork.Client from ClientMemMapStorage
	m.Rcv <- func() error {
		for _, c := range ads {
			delete(m.storage, c.ID.String())
		}
		return nil
	}
	return <-m.Err
}

func (m *MemMapStorage) cleanup() {
	close(m.Rcv)
	close(m.Err)
	m.storage = nil
	m.logger.Info().Msg("Cleanup Done.")
}
