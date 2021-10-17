package usecase

import (
	"context"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
)

// --------------------- IsAdmin ------------------------

type IsAdminCmd func(ctx context.Context, input IsAdminInput) (bool, error)

// IsAdminInput .
type IsAdminInput struct {
	AccountID domain.AccountID `json:"account_id,required"`
}

// IsAdmin is the UseCase handler that create a domain.Account.
func (u UseCase) IsAdmin() IsAdminCmd {
	return func(ctx context.Context, input IsAdminInput) (bool, error) {
		acc, err := u.DB.ByID(ctx, input.AccountID)
		if err != nil {
			return false, err
		}
		return acc.IsAdmin(), err
	}
}
