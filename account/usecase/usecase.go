package usecase

import (
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Database interface {
	Create(...*domain.Account) error
	Update(...*domain.Account) error
	All() ([]*domain.Account, error)
	ByID(id domain.AccountID) (*domain.Account, error)
	ByEmail(email string) ([]*domain.Account, error)
	ByFirstname(firstname string) ([]*domain.Account, error)
	ByLastname(lastname string) ([]*domain.Account, error)
	ByFullname(firstname, lastname string) ([]*domain.Account, error)
	Remove(...*domain.Account) error
}

type UseCase struct {
	DB     Database
	logger zerolog.Logger
}

func NewUseCase(db Database) *UseCase {
	return &UseCase{
		DB:     db,
		logger: log.With().Str("service", "UseCase").Logger(),
	}
}
