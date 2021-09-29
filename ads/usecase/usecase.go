package usecase

import (
	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type database interface {
	All() ([]*domain.Ad, error)
	ByID(id domain.AdID) (*domain.Ad, error)
	Save(...*domain.Ad) error
	Remove(...*domain.Ad) error
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
