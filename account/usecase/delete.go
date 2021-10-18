package usecase

import (
	"context"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
)

// DeleteAccountCmd is a func type that return the logic for deleting a domain.Account.
type DeleteAccountCmd func(ctx context.Context, input DeleteAccountInput) (*domain.Account, error)

// DeleteAccountInput is used by UseCase.DeleteAccount for the deletion of an account.
type DeleteAccountInput struct {
	ID       domain.AccountID `json:"id" binding:"required"`
	Email    string           `json:"email" binding:"required"`
	Password string           `json:"password" binding:"required"`
}

// DeleteAccount is the UseCase handler that delete a domain.Account.
func (u UseCase) DeleteAccount() DeleteAccountCmd {
	return func(ctx context.Context, input DeleteAccountInput) (*domain.Account, error) {
		acc, err := u.DB.Remove(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		if err = u.AuthService.Unregister(acc.Email, input.Email, input.ID); err != nil {
			return nil, err
		}
		if err = u.BalanceService.Delete(input.ID); err != nil {
			return nil, err
		}
		return acc, err
	}
}
