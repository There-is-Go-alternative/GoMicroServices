package usecase

import (
	"context"
	"fmt"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/xerrors"
	"time"
)

// --------------------- CreateAccount ------------------------

// CreateAccountCmd is a func type that return the logic for creating a domain.Account.
type CreateAccountCmd func(ctx context.Context, input CreateAccountInput) (*domain.Account, error)

// CreateAccountInput is used by UseCase.CreateAccount for the creation of an account.
type CreateAccountInput struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// CreateAccount is the UseCase handler that create a domain.Account.
func (u UseCase) CreateAccount() CreateAccountCmd {
	return func(ctx context.Context, input CreateAccountInput) (*domain.Account, error) {
		// Dealing with CreateAccountInput
		account := &domain.Account{Email: input.Email, Firstname: input.Firstname, Lastname: input.Lastname}
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

		// Adding createdAt timestamp
		account.CreatedAt = time.Now()
		if err = u.AuthService.Register(input.Email, input.Password, account.ID); err != nil {
			return nil, fmt.Errorf("error when calling Auth service Create: %v", err)
		}
		if err = u.BalanceService.Create(account.ID); err != nil {
			return nil, fmt.Errorf("error when calling Auth service Create: %v", err)
		}

		// Creating account
		account, err = u.DB.Create(ctx, account)
		if err != nil {
			return nil, err
		}
		return account, nil
	}
}
