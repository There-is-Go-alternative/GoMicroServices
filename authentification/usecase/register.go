package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
)

type RegisterProto func(ctx context.Context, input domain.Register) (string, error)

func (u UseCase) Register() RegisterProto {
	return func(ctx context.Context, input domain.Register) (string, error) {

		_, err := u.DB.Save(ctx, input)
		if err != nil {
			return "", err
		}

		return "success", nil
	}
}
