package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/There-is-Go-alternative/GoMicroServices/funds/domain"
)

type CreateFundsCmd func(ctx context.Context, input CreateFundsInput) (*domain.Funds, error)

type CreateFundsInput struct {
	UserId  string `json:"user_id" binding:"required"`
	Balance int    `json:"initial_balance"`
}

func (u UseCase) Create() CreateFundsCmd {
	return func(ctx context.Context, input CreateFundsInput) (*domain.Funds, error) {
		if id, err := domain.NewFundsID(); err != nil {
			return nil, err
		} else {
			funds := &domain.Funds{
				ID:          *id,
				UserId:      input.UserId,
				Balance:     input.Balance,
				LastUpdated: time.Now(),
			}

			_, getErr := u.DB.ByUserId(ctx, input.UserId)

			if getErr == nil {
				return nil, fmt.Errorf("this user (%v) already have a balance", input.UserId)
			}

			createErr := u.DB.Create(ctx, *funds)

			fmt.Println("Before return")

			if createErr != nil {
				fmt.Println("Error there is %v", createErr)
			}

			return funds, createErr
		}
	}
}
