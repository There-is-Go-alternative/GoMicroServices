package usecase

import (
	"context"
	"fmt"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/xerrors"
)

type CreateAccountCmd func(ctx context.Context, input CreateAccountInput) (*domain.Account, error)

type CreateAccountInput struct {
	Email string `json:"email" binding:"required"`
}

func (u UseCase) CreateAccount() CreateAccountCmd {
	return func(ctx context.Context, input CreateAccountInput) (*domain.Account, error) {
		// Dealing with CreateAccountInput
		account := &domain.Account{Email: input.Email}
		accountID, err := domain.NewAccountID()
		if err != nil {
			return nil, err
		}
		account.ID = *accountID

		// Validating sent data
		if err = account.Validate(); err != nil {
			return nil, xerrors.ErrorWithCode{
				Code: xerrors.CodeInvalidData, Err: fmt.Errorf("invalid user account data: %v", err),
			}
		}
		// Checking in Database if there is no account with this ID.
		duplicate, err := u.DB.ByID(ctx, *accountID)
		if err == nil {
			return nil, fmt.Errorf("account ID (%v) already exist in DB: %v", accountID, duplicate)
		}

		// Creating account
		err = u.DB.Create(ctx, account)
		return account, err
	}
}
