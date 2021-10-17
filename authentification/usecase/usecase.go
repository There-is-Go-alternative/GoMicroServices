package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

// Database is an interface that represent all possible actions that can be performed on a domain.Account DB.

type Database interface {
	FindByEmail(context.Context, string) (domain.Auth, error)
	FindByID(context.Context, string) (domain.Auth, error)
	Save(context.Context, domain.Auth) (*mongo.InsertOneResult, error)
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
