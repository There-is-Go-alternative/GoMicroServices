package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/funds/domain"
)

type AllCmd func(ctx context.Context) ([]*domain.Funds, error)

func (u UseCase) All() AllCmd {
	return func(ctx context.Context) ([]*domain.Funds, error) {
		return u.DB.All(ctx)
	}
}

type GetByUserIDCmd func(ctx context.Context, id string) (*domain.Funds, error)

func (u UseCase) GetByUserID() GetByUserIDCmd {
	return func(ctx context.Context, id string) (*domain.Funds, error) {
		return u.DB.ByUserId(ctx, id)
	}
}

type GetByIDCmd func(ctx context.Context, id domain.FundsID) (*domain.Funds, error)

func (u UseCase) GetByID() GetByIDCmd {
	return func(ctx context.Context, id domain.FundsID) (*domain.Funds, error) {
		return u.DB.ById(ctx, &id)
	}
}
