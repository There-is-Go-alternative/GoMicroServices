package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/funds/domain"
)

type DeleteByUserIDCmd func(ctx context.Context, id string) error

func (u UseCase) DeleteByUserID() DeleteByUserIDCmd {
	return func(ctx context.Context, id string) error {
		return u.DB.DeleteByUserId(ctx, id)
	}
}

type DeleteByIDCmd func(ctx context.Context, id domain.FundsID) error

func (u UseCase) DeleteByID() DeleteByIDCmd {
	return func(ctx context.Context, id domain.FundsID) error {
		return u.DB.DeleteById(ctx, &id)
	}
}
