package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
)

type DeleteAccountProto func(ctx context.Context, input domain.Auth) error

func (u UseCase) DeleteAccount() DeleteAccountProto {
	return func(ctx context.Context, input domain.Auth) error {

		err := u.DB.Delete(ctx, input.ID)
		if err != nil {
			return err
		}

		return nil
	}
}
