package usecase

import (
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type database interface {
	All() ([]*domain.Account, error)
	ByID(id domain.AccountID) (*domain.Account, error)
	Save(...*domain.Account) error
	Remove(...*domain.Account) error
}

type UseCase struct {
	DB     database
	logger zerolog.Logger
}

func NewGetUseCase(db database) *UseCase {
	return &UseCase{
		DB:     db,
		logger: log.With().Str("service", "UseCase").Logger(),
	}
}
