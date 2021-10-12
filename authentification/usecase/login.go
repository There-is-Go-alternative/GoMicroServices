package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
)

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginProto func(ctx context.Context, input LoginDTO) (*domain.Token, error)

func (u UseCase) Login() LoginProto {
	return func(ctx context.Context, input LoginDTO) (*domain.Token, error) {
		var auth domain.Auth
		auth, err := u.DB.FindByEmail(ctx, input.Email)
		if err != nil {
			return nil, err
		}

		hashed, err := domain.HashPassword(input.Password)
		if err != nil {
			return nil, err
		}
		err = domain.VerifyPassword(hashed, auth.Password)
		if err != nil {
			return nil, err
		}

		token, err := domain.CreateToken(auth.ID)
		if err != nil {
			return nil, err
		}

		return &domain.Token{
			Token: token,
		}, nil
	}
}
