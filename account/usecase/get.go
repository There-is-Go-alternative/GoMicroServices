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

//type GetAccountByIdCmd func(ctx context.Context, accountId string, token string) (*domain.PartialAccount, error)
//
//func GetAccountById(authClient authPb.AuthServiceClient) GetAccountByIdCmd {
//	return func(ctx context.Context, accountId string, token string) (*domain.PartialAccount, error) {
//		_, e := authClient.Authorize(context.Background(), &authPb.AuthorizationRequest{
//			Token: token,
//		})
//		if e != nil {
//			log.Print("/GET account/:id: Error while authenticating")
//			return nil, e
//		}
//
//		var account domain.Account
//		res := database.DB.Where("id = ?", accountId).First(&account).Error
//
//		return &domain.PartialAccount{
//			Email: account.Email,
//			Lastname: account.Lastname,
//		}, res
//	}
//}
//
//type GetMeCmd func(ctx context.Context, token string) (*domain.Account, error)
//
//func GetMe(authClient authPb.AuthServiceClient) GetMeCmd {
//	return func(ctx context.Context, token string) (*domain.Account, error) {
//		authRes, e := authClient.Authorize(context.Background(), &authPb.AuthorizationRequest{
//			Token: token,
//		})
//		if e != nil {
//			log.Print("/GET me: Error while authenticating")
//			return nil, e
//		}
//
//		var account domain.Account
//		err := database.DB.Where("id = ?", authRes.UserId).First(&account).Error
//
//		return &account, err
//	}
//}
