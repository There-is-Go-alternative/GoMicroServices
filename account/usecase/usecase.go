package usecase

import (
	"context"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Database interface {
	Create(ctx context.Context, accounts ...*domain.Account) error
	Update(ctx context.Context, accounts ...*domain.Account) error
	All(ctx context.Context) ([]*domain.Account, error)
	ByID(ctx context.Context, id domain.AccountID) (*domain.Account, error)
	ByEmail(ctx context.Context, email string) ([]*domain.Account, error)
	ByFirstname(ctx context.Context, firstname string) ([]*domain.Account, error)
	ByLastname(ctx context.Context, lastname string) ([]*domain.Account, error)
	ByFullname(ctx context.Context, firstname, lastname string) ([]*domain.Account, error)
	Remove(ctx context.Context, accounts ...*domain.Account) error
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
