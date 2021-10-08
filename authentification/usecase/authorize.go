package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
)

type AuthorizeProto func(ctx context.Context, token domain.Token) (*domain.Authorize, error)

func (u UseCase) Authorize() AuthorizeProto {
	return func(ctx context.Context, token domain.Token) (*domain.Authorize, error) {
		id, err := domain.VerifyToken(token.Token)
		if err != nil {
			return nil, err
		}

		return &domain.Authorize{
			UserID: id,
		}, nil
	}
}
