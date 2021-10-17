package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal"
)

type BuyAdCmd func(ctx context.Context, id domain.AdID) (*domain.Ad, error)

func (u UseCase) BuyAd() BuyAdCmd {
	return func(ctx context.Context, id domain.AdID) (*domain.Ad, error) {
		ad, err := u.DB.ByID(ctx, id)

		if err != nil {
			return nil, internal.NewInternalError(internal.NotFound, internal.AdNotFound)
		}

		if ad.State == "close" {
			return nil, internal.NewInternalError(internal.BadRequest, internal.AdIsClose)
		}

		ad.State = "close"

		err = u.DB.Update(ctx, ad)

		if err != nil {
			return nil, internal.NewInternalError(internal.DatabaseError, internal.InternalServerError)
		}

		return ad, nil
	}
}
