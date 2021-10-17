package usecase

import (
	"context"
	"time"

	"github.com/There-is-Go-alternative/GoMicroServices/transactions/domain"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Database interface {
	All(ctx context.Context) ([]*domain.Transaction, error)
	ById(ctx context.Context, id string) (*domain.Transaction, error)
	ByDateRange(ctx context.Context, start time.Time, end time.Time) ([]*domain.Transaction, error)
	BySellerId(ctx context.Context, id string) ([]*domain.Transaction, error)
	ByBuyerId(ctx context.Context, id string) ([]*domain.Transaction, error)
	ByBuyerIdDateRange(ctx context.Context, id string, start time.Time, end time.Time) ([]*domain.Transaction, error)
	BySellerIdDateRange(ctx context.Context, id string, start time.Time, end time.Time) ([]*domain.Transaction, error)
	Register(ctx context.Context, t *domain.Transaction) error
}

type AccountService interface {
	GetUser(ctx context.Context, token string) (*domain.Account, error)
	IsAdmin(id string) (bool, error)
}

type FundsService interface {
	Increase(ctx context.Context, user_id string, by float64) error
	Decrease(ctx context.Context, user_id string, by float64) error
	GetBalance(ctx context.Context, user_id string) (*float64, error)
}

type AdsService interface {
	GetAd(ctx context.Context, token string) (*domain.Ad, error)
	Buy(ctx context.Context, id, token string) error
}

type UseCase struct {
	DB     Database
	logger zerolog.Logger
	Auth   AccountService
	Ads    AdsService
	Funds  FundsService
}

func NewUseCase(auth AccountService, ads AdsService, funds FundsService, db Database) *UseCase {
	return &UseCase{
		DB:     db,
		logger: log.With().Str("service", "UseCase").Logger(),
		Auth:   auth,
		Ads:    ads,
		Funds:  funds,
	}
}
