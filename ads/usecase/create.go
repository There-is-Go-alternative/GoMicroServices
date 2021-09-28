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
		account := &domain.Account{Email: input.Email}
		accountID, err := domain.NewAccountID()
		if err != nil {
			// TODO: BETTER ERROR
			return nil, err
		}
		account.ID = accountID
		if !account.Validate() {
			return nil, xerrors.ErrorWithCode{
				Code: xerrors.CodeInvalidData, Err: fmt.Errorf("invalid user account data: %v", account),
			}
		}
		err = u.DB.Save(account)
		return account, err
	}
}
