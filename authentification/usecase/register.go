package usecase

import (
	"context"
	"log"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
)

type RegisterProto func(ctx context.Context, input domain.Auth) error

func (u UseCase) Register() RegisterProto {
	return func(ctx context.Context, input domain.Auth) error {

		err := u.DB.Save(ctx, input)
		if err != nil {
			log.Println(err)
			return err
		}

		return nil
	}
}
