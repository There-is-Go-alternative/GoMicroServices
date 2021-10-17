package usecase

import (
	"context"
	"fmt"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/xerrors"
)

// --------------------- GetAllAccounts ------------------------

// GetAllAccountsCmd is a func type that return the logic for retrieving all domain.Account.
type GetAllAccountsCmd func(ctx context.Context) ([]*domain.Account, error)

// GetAllAccounts is the UseCase handler that retrieve all domain.Account.
func (u UseCase) GetAllAccounts() GetAllAccountsCmd {
	return func(ctx context.Context) ([]*domain.Account, error) {
		u.logger.Info("Fetching all accounts ...")
		defer u.logger.Info("All accounts fetched !")
		accounts, err := u.DB.All(ctx)
		if err != nil {
			return nil, err
		}
		balances, err := u.BalanceService.GetAll()
		if err != nil {
			return nil, fmt.Errorf("error when calling balance service: %v", err)
		}
		for _, account := range accounts {
			b, ok := balances[account.ID]
			if ok {
				account.Balance = b.Balance
			}
		}
		return accounts, nil
	}
}

// --------------------- GetAccountByID ------------------------

// GetAccountByIDCmd is a func type that return the logic for retrieving a domain.Account by a specified ID.
type GetAccountByIDCmd func(ctx context.Context, id domain.AccountID) (*domain.Account, error)

// GetAccountByIDInput is used by UseCase.GetAccountByID for the retrieval of an account.
type GetAccountByIDInput struct {
	Email string `json:"id" binding:"required"`
}

// GetAccountByID is the UseCase handler that retrieve a domain.Account by a domain.AccountID.
func (u UseCase) GetAccountByID() GetAccountByIDCmd {
	return func(ctx context.Context, id domain.AccountID) (*domain.Account, error) {
		// TODO: Add auth service check here
		u.logger.Infof("Fetching account by id: %v", id)
		defer u.logger.Infof("account %v fetched !", id)
		account, err := u.DB.ByID(ctx, id)
		if err != nil {
			return nil, err
		}
		balance, err := u.BalanceService.GetByID(account.ID)
		if err != nil {
			return nil, err
		}
		account.Balance = *balance
		return account, nil
	}
}

// --------------------- GetAccountByEmail ------------------------

// GetAccountByEmailCmd is a func type that return the logic for retrieving the Account by a specified Email
// Strict: As the Email is unique by account, only one domain.Account will be fetched.
type GetAccountByEmailCmd func(ctx context.Context, email EmailInputCmd) (*domain.Account, error)

// EmailInputCmd is used by UseCase.GetAccountByEmail for the retrieval of an account.
type EmailInputCmd struct {
	Email string `json:"email" binding:"required"`
}

// GetAccountByEmail is the UseCase handler that retrieve a domain.Account by an Email.
func (u UseCase) GetAccountByEmail() GetAccountByEmailCmd {
	return func(ctx context.Context, cmd EmailInputCmd) (*domain.Account, error) {
		// TODO: Add auth service check here
		u.logger.Infof("Fetching account by email: %v", cmd.Email)
		defer u.logger.Info("AccountByEmail fetched !")

		// Fetching account by Firstname
		accounts, err := u.DB.ByEmail(ctx, cmd.Email)
		if err != nil {
			return nil, err
		}

		// If several account with email found with same Firstname, Database dirty
		if len(accounts) > 1 {
			return nil, fmt.Errorf("serveral account found with email: %v", cmd.Email)
		}

		// If not account found by email, user doesn't exist
		if len(accounts) == 0 {
			return nil, xerrors.AccountNotFound
		}

		// return first found user
		return accounts[0], nil
	}
}

// --------------------- GetAccountByFirstname ------------------------

// GetAccountByFirstnameCmd is a func type that return the logic for retrieving all domain.Account by a Firstname
type GetAccountByFirstnameCmd func(context.Context, FirstnameInputCmd) ([]*domain.Account, error)

// FirstnameInputCmd is used by UseCase.GetAccountByFirstname for the retrieval of an account.
type FirstnameInputCmd struct {
	Firstname string `json:"firstname" binding:"required"`
}

// GetAccountByFirstname is the UseCase handler that retrieve all domain.Account by a Firstname.
func (u UseCase) GetAccountByFirstname() GetAccountByFirstnameCmd {
	return func(ctx context.Context, cmd FirstnameInputCmd) ([]*domain.Account, error) {
		// TODO: Add auth service check here
		u.logger.Infof("Fetching account by Firstname: %v", cmd.Firstname)
		defer u.logger.Info("All GetAccountByFirstname fetched !")

		// Fetching account by Firstname
		return u.DB.ByFirstname(ctx, cmd.Firstname)
	}
}

// --------------------- GetAccountByLastname ------------------------

// GetAccountByLastnameCmd is a func type that return the logic for retrieving all domain.Account by a Lastname
type GetAccountByLastnameCmd func(context.Context, LastnameInputCmd) ([]*domain.Account, error)

// LastnameInputCmd is used by UseCase.GetAccountByLastname for the retrieval of an account.
type LastnameInputCmd struct {
	Lastname string `json:"lastname" binding:"required"`
}

// GetAccountByLastname is the UseCase handler that retrieve all domain.Account by a Lastname.
func (u UseCase) GetAccountByLastname() GetAccountByLastnameCmd {
	return func(ctx context.Context, cmd LastnameInputCmd) ([]*domain.Account, error) {
		// TODO: Add auth service check here
		u.logger.Infof("Fetching account by Lastname: %v", cmd.Lastname)
		defer u.logger.Info("All AccountByLastname fetched !")

		// Fetching account by Lastname
		return u.DB.ByLastname(ctx, cmd.Lastname)
	}
}

// --------------------- GetAccountByFullname ------------------------

// GetAccountByFullnameCmd is a func type that return the logic for retrieving all domain.Account by a Fullname.
type GetAccountByFullnameCmd func(context.Context, FullnameInputCmd) ([]*domain.Account, error)

// FullnameInputCmd is used by UseCase.GetAccountByFullname for the retrieval of a domain.Account.
type FullnameInputCmd struct {
	FirstnameInputCmd
	LastnameInputCmd
}

// GetAccountByFullname is the UseCase handler that retrieve all domain.Account by a Fullname.
func (u UseCase) GetAccountByFullname() GetAccountByFullnameCmd {
	return func(ctx context.Context, cmd FullnameInputCmd) ([]*domain.Account, error) {
		// TODO: Add auth service check here
		u.logger.Infof("Fetching account by Fullname: %v %v", cmd.Firstname, cmd.Lastname)
		defer u.logger.Info("All accounts fetched !")

		// Fetching account by Fullname
		return u.DB.ByFullname(ctx, cmd.Firstname, cmd.Lastname)
	}
}
