package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/transactions/domain"
)

type GetByIdCmd func(ctx context.Context, id string) (*domain.Transaction, error)

func (u UseCase) GetById() GetByIdCmd {
	return func(ctx context.Context, id string) (*domain.Transaction, error) {
		return u.DB.ById(ctx, id)
	}
}

type GetAllCmd func(ctx context.Context) ([]*domain.Transaction, error)

func (u UseCase) GetAll() GetAllCmd {
	return func(ctx context.Context) ([]*domain.Transaction, error) {
		return u.DB.All(ctx)
	}
}

type GetBySellerIdCmd func(ctx context.Context, id string) ([]*domain.Transaction, error)

func (u UseCase) GetBySellerId() GetBySellerIdCmd {
	return func(ctx context.Context, id string) ([]*domain.Transaction, error) {
		return u.DB.BySellerId(ctx, id)
	}
}

type GetByBuyerIdCmd func(ctx context.Context, id string) ([]*domain.Transaction, error)

func (u UseCase) GetByBuyerId() GetByBuyerIdCmd {
	return func(ctx context.Context, id string) ([]*domain.Transaction, error) {
		return u.DB.ByBuyerId(ctx, id)
	}
}
