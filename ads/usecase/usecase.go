package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type database interface {
	All(ctx context.Context) ([]*domain.Ad, error)
	ByID(ctx context.Context, id domain.AdID) (*domain.Ad, error)
	Create(ctx context.Context, ads ...*domain.Ad) error
	Remove(ctx context.Context, ads ...*domain.Ad) error
	Update(ctx context.Context, ad ...*domain.Ad) error
	Search(ctx context.Context, content string) ([]domain.Ad, error)
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
