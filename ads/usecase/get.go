package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
)

type GetAllAdsCmd func(ctx context.Context) ([]*domain.Ad, error)

func (u UseCase) GetAllAds() GetAllAdsCmd {
	return func(ctx context.Context) ([]*domain.Ad, error) {
		u.logger.Info().Msg("Fetching all ads ...")
		defer u.logger.Info().Msg("All ads fetched !")
		//TODO: In the future, when the Database will be real, find another way to make this code return an empty array
		all_ads, err := u.DB.All(ctx)

		if all_ads == nil {
			return make([]*domain.Ad, 0), nil
		}
		return all_ads, err
	}
}

type GetAdByIdCmd func(ctx context.Context, id domain.AdID) (*domain.Ad, error)

func (u UseCase) GetAdById() GetAdByIdCmd {
	return func(ctx context.Context, id domain.AdID) (*domain.Ad, error) {
		u.logger.Info().Msgf("Fetching ad by id: %v", id)
		defer u.logger.Info().Msg("All ads fetched !")
		return u.DB.ByID(ctx, id)
	}
}
