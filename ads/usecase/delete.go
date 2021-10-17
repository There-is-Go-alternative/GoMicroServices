package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
)

type DeleteAdCmd func(ctx context.Context, input DeleteAdInput) (*domain.Ad, error)

type DeleteAdInput struct {
	ID domain.AdID `json:"id" binding:"required"`
}

func (u UseCase) DeleteAd() DeleteAdCmd {
	return func(ctx context.Context, input DeleteAdInput) (*domain.Ad, error) {
		ad, err := u.DB.ByID(ctx, input.ID)

		if err != nil {
			return nil, nil
		}
		err = u.DB.Remove(ctx, ad)
		return ad, err
	}
}
