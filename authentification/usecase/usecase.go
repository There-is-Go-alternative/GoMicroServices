package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
)

// Database is an interface that represent all possible actions that can be performed on a domain.Account DB.

type AuthUseCase interface {
	Login(ctx context.Context, input LoginDTO) (*domain.Token, error)
}
