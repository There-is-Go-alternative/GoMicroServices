package usecase

import (
	"context"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
)

type GetAllAccountsCmd func(ctx context.Context) ([]*domain.Account, error)

func (u UseCase) GetAllAccounts() GetAllAccountsCmd {
	return func(ctx context.Context) ([]*domain.Account, error) {
		u.logger.Info().Msg("Fetching all accounts ...")
		defer u.logger.Info().Msg("All accounts fetched !")
		return u.DB.All()
	}
}

type GetAccountByIdCmd func(ctx context.Context, id domain.AccountID) (*domain.Account, error)

func (u UseCase) GetAccountById() GetAccountByIdCmd {
	return func(ctx context.Context, id domain.AccountID) (*domain.Account, error) {
		// TODO: Add auth service check here
		u.logger.Info().Msgf("Fetching account by id: %v", id)
		defer u.logger.Info().Msg("All accounts fetched !")
		return u.DB.ByID(id)
	}
}
