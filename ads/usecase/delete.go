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
		ad := domain.Ad{ID: input.ID}
		err := u.DB.Remove(&ad)
		return &ad, err
	}
}
