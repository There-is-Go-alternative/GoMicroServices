package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal"
)

type CreateAdCmd func(ctx context.Context, input CreateAdInput) (*domain.Ad, error)

type CreateAdInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       *uint   `json:"price" binding:"required"`
	Pictures   	[]string `json:"pictures" binding:"required"`
	UserId		string `json:"userid"`
}

func (u UseCase) CreateAd() CreateAdCmd {
	return func(ctx context.Context, input CreateAdInput) (*domain.Ad, error) {
		ad := &domain.Ad{Title: input.Title, Description: input.Description, Price: *input.Price, Pictures: input.Pictures, UserId: input.UserId, State: "open"}
		AdID, err := domain.NewAdID()

		if err != nil {
			return nil, internal.NewInternalError(internal.InternalServerError, internal.InternalServerErrorMsg)
		}
		ad.ID = AdID
		if !ad.Validate() {
			return nil, internal.NewInternalError(internal.BadRequest, internal.BadRequestMsg)
		}
		err = u.DB.Create(ctx, ad)

		if err != nil {
			return nil, internal.NewInternalError(internal.InternalServerError, internal.InternalServerErrorMsg)
		}
		return ad, nil
	}
}
