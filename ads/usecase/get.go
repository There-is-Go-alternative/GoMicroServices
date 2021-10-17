package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal"
)

type GetAllAdsCmd func(ctx context.Context) ([]*domain.Ad, error)

func (u UseCase) GetAllAds() GetAllAdsCmd {
	return func(ctx context.Context) ([]*domain.Ad, error) {
		u.logger.Info().Msg("Fetching all ads ...")
		defer u.logger.Info().Msg("All ads fetched !")
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
		ads, err := u.DB.ByID(ctx, id)

		if err != nil {
			return nil, internal.NewInternalError(internal.InternalServerError, internal.InternalServerErrorMsg)
		}

		return ads, nil
	}
}
