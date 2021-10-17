package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal"
)

type DeleteAdCmd func(ctx context.Context, input DeleteAdInput) (*domain.Ad, error)

type DeleteAdInput struct {
	ID domain.AdID `json:"id" binding:"required"`
}

func (u UseCase) DeleteAd() DeleteAdCmd {
	return func(ctx context.Context, input DeleteAdInput) (*domain.Ad, error) {
		ad, err := u.DB.ByID(ctx, input.ID)

		if err != nil {
			return nil, internal.NewInternalError(internal.InternalServerError, internal.InternalServerErrorMsg)
		}
		err = u.DB.Remove(ctx, ad)

		if err != nil {
			return nil, internal.NewInternalError(internal.InternalServerError, internal.InternalServerErrorMsg)
		}
		return ad, nil
	}
}
