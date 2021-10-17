package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/funds/domain"
)

type IncreaseDecreaseInput struct {
	By float64 `json:"by" binding:"required"`
}

type SetInput struct {
	NewBalance float64 `json:"new_balance" binding:"required"`
}

type IncreaseByUserCmd func(ctx context.Context, id string, input IncreaseDecreaseInput) error

func (u UseCase) IncreaseByUser() IncreaseByUserCmd {
	return func(ctx context.Context, id string, input IncreaseDecreaseInput) error {
		return u.DB.IncreaseByUser(ctx, id, input.By)
	}
}

type DecreaseByUserCmd func(ctx context.Context, id string, input IncreaseDecreaseInput) error

func (u UseCase) DecreaseByUser() DecreaseByUserCmd {
	return func(ctx context.Context, id string, input IncreaseDecreaseInput) error {
		return u.DB.DecreaseByUser(ctx, id, input.By)
	}
}

type SetByUserCmd func(ctx context.Context, id string, input SetInput) error

func (u UseCase) SetByUser() SetByUserCmd {
	return func(ctx context.Context, id string, input SetInput) error {
		return u.DB.UpdateByUser(ctx, id, input.NewBalance)
	}
}

type IncreaseCmd func(ctx context.Context, id domain.FundsID, input IncreaseDecreaseInput) error

func (u UseCase) Increase() IncreaseCmd {
	return func(ctx context.Context, id domain.FundsID, input IncreaseDecreaseInput) error {
		return u.DB.Increase(ctx, &id, input.By)
	}
}

type DecreaseCmd func(ctx context.Context, id domain.FundsID, input IncreaseDecreaseInput) error

func (u UseCase) Decrease() DecreaseCmd {
	return func(ctx context.Context, id domain.FundsID, input IncreaseDecreaseInput) error {
		return u.DB.Decrease(ctx, &id, input.By)
	}
}

type SetCmd func(ctx context.Context, id domain.FundsID, input SetInput) error

func (u UseCase) Set() SetCmd {
	return func(ctx context.Context, id domain.FundsID, input SetInput) error {
		return u.DB.Update(ctx, &id, input.NewBalance)
	}
}
