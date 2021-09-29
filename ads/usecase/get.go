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
		return u.DB.All()
	}
}

type GetAdByIdCmd func(ctx context.Context, id domain.AdID) (*domain.Ad, error)

func (u UseCase) GetAdById() GetAdByIdCmd {
	return func(ctx context.Context, id domain.AdID) (*domain.Ad, error) {
		// TODO: Add auth service check here
		u.logger.Info().Msgf("Fetching ad by id: %v", id)
		defer u.logger.Info().Msg("All ads fetched !")
		return u.DB.ByID(id)
	}
}
