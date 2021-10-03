package usecase

import (
	"context"
	"fmt"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal/xerrors"
)

type CreateAdCmd func(ctx context.Context, input CreateAdInput) (*domain.Ad, error)

type CreateAdInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       uint   `json:"price" binding:"required"`
	Picture     string `json:"picture" binding:"required"`
}

func (u UseCase) CreateAd() CreateAdCmd {
	return func(ctx context.Context, input CreateAdInput) (*domain.Ad, error) {
		ad := &domain.Ad{Title: input.Title, Description: input.Description, Price: input.Price, Picture: input.Picture}
		AdID, err := domain.NewAdID()
		if err != nil {
			// TODO: BETTER ERROR
			return nil, err
		}
		ad.ID = AdID
		if !ad.Validate() {
			return nil, xerrors.ErrorWithCode{
				Code: xerrors.CodeInvalidData, Err: fmt.Errorf("invalid user ad data: %v", ad),
			}
		}
		err = u.DB.Create(ctx, ad)
		return ad, err
	}
}
