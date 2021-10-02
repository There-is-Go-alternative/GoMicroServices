package usecase

import (
	"context"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
)

type DeleteAccountCmd func(ctx context.Context, input DeleteAccountInput) (*domain.Account, error)

type DeleteAccountInput struct {
	ID domain.AccountID `json:"id" binding:"required"`
}

func (u UseCase) DeleteAccount() DeleteAccountCmd {
	return func(ctx context.Context, input DeleteAccountInput) (*domain.Account, error) {
		account := domain.Account{ID: input.ID}
		err := u.DB.Remove(ctx, &account)
		return &account, err
	}
}
