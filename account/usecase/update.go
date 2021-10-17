package usecase

import (
	"context"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/xerrors"
	"github.com/imdario/mergo"
	"time"
)

// --------------------- UpdateAccount ------------------------

// UpdateAccountCmd is a func type that return the logic for updating a domain.Account.
type UpdateAccountCmd func(ctx context.Context, input UpdateAccountInput) (*domain.Account, error)

// UpdateAccountInput is used by UseCase.UpdateAccount for the update of a domain.Account.
type UpdateAccountInput struct {
	domain.Account
}

// UpdateAccount is the UseCase handler that retrieve all domain.Account by a Fullname.
func (u UseCase) UpdateAccount() UpdateAccountCmd {
	return func(ctx context.Context, cmd UpdateAccountInput) (*domain.Account, error) {
		if err := cmd.Account.Validate(); err != nil {
			return nil, err
		}
		a, err := u.DB.ByID(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		err = mergo.Merge(a, &cmd.Account, mergo.WithOverride)
		if err != nil {
			return nil, xerrors.ErrorWithCode{Code: xerrors.InternalError, Err: err}
		}
		a.UpdatedAt = time.Now()
		a, err = u.DB.Update(ctx, a)
		if err != nil {
			return nil, err
		}
		return a, nil
	}
}

// --------------------- PatchAccount ------------------------

// PatchAccountCmd is a func type that return the logic for patching a domain.Account.
type PatchAccountCmd func(ctx context.Context, input PatchAccountInput) (*domain.Account, error)

// PatchAccountInput is used by UseCase.UpdateAccount for the patch of a domain.Account.
type PatchAccountInput struct {
	domain.Account
}

// PatchAccount is the UseCase handler that retrieve all domain.Account by a Fullname.
func (u UseCase) PatchAccount() PatchAccountCmd {
	return func(ctx context.Context, cmd PatchAccountInput) (*domain.Account, error) {
		a, err := u.DB.ByID(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		// Merging Structs
		err = mergo.Merge(a, &cmd.Account, mergo.WithOverride)
		if err != nil {
			return nil, xerrors.ErrorWithCode{Code: xerrors.InternalError, Err: err}
		}
		// TODO: Check Unique fields ??

		// Updating modification timestamp
		a.UpdatedAt = time.Now()

		// Updating account in DB
		if a, err = u.DB.Update(ctx, a); err != nil {
			return nil, err
		}
		return a, nil
	}
}
