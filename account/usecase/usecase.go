package usecase

import (
	"context"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	log "github.com/sirupsen/logrus"
)

// Database is an interface that represent all possible actions that can be performed on a domain.Account DB.
type Database interface {
	Create(ctx context.Context, account *domain.Account) (*domain.Account, error)
	Update(ctx context.Context, account *domain.Account) (*domain.Account, error)
	All(ctx context.Context) ([]*domain.Account, error)
	ByID(ctx context.Context, id domain.AccountID) (*domain.Account, error)
	ByEmail(ctx context.Context, email string) ([]*domain.Account, error)
	ByFirstname(ctx context.Context, firstname string) ([]*domain.Account, error)
	ByLastname(ctx context.Context, lastname string) ([]*domain.Account, error)
	ByFullname(ctx context.Context, firstname, lastname string) ([]*domain.Account, error)
	Remove(ctx context.Context, id domain.AccountID) (*domain.Account, error)
}

type AuthService interface {
	Authorize(token string) (domain.AccountID, error)
	Register(email, password string) error
	Unregister(email string) error
}

type BalanceService interface {
	Authorize(token string) (domain.AccountID, error)
}

// UseCase handle the business logic
type UseCase struct {
	AuthService    AuthService
	BalanceService BalanceService
	DB             Database
	logger         *log.Logger
}

// NewUseCase return an initialized UseCase, using Database
func NewUseCase(auth AuthService, balance BalanceService, db Database, logger *log.Logger) *UseCase {
	return &UseCase{
		AuthService:    auth,
		BalanceService: balance,
		DB:             db,
		logger:         logger.WithField("service", "UseCase").Logger,
	}
}
