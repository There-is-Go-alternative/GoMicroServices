package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/funds/domain"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Database interface {
	ByUserId(ctx context.Context, id string) (*domain.Funds, error)
	ById(ctx context.Context, id *domain.FundsID) (*domain.Funds, error)
	Update(ctx context.Context, id *domain.FundsID, new_balance float64) error
	UpdateByUser(ctx context.Context, id string, new_balance float64) error
	Increase(ctx context.Context, id *domain.FundsID, by float64) error
	IncreaseByUser(ctx context.Context, id string, by float64) error
	Decrease(ctx context.Context, id *domain.FundsID, by float64) error
	DecreaseByUser(ctx context.Context, id string, by float64) error
	Create(ctx context.Context, f domain.Funds) error
	DeleteByUserId(ctx context.Context, id string) error
	DeleteById(ctx context.Context, id *domain.FundsID) error
	All(ctx context.Context) ([]*domain.Funds, error)
}

type AuthService interface {
	ValidateToken(token string) error
}

type UseCase struct {
	DB     Database
	logger zerolog.Logger
	Auth   AuthService
}

func NewUseCase(auth AuthService, db Database) *UseCase {
	return &UseCase{
		DB:     db,
		logger: log.With().Str("service", "UseCase").Logger(),
		Auth:   auth,
	}
}
