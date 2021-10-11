package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal/xerrors"
)

type CreateAdCmd func(ctx context.Context, input CreateAdInput) (*domain.Ad, error)

type CreateAdInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       *uint   `json:"price" binding:"required"`
	Pictures   	string `json:"pictures" binding:"required"`
	UserId		string `json:"userid"`
}

func (u UseCase) CreateAd() CreateAdCmd {
	return func(ctx context.Context, input CreateAdInput) (*domain.Ad, error) {
		ad := &domain.Ad{Title: input.Title, Description: input.Description, Price: *input.Price, Pictures: strings.Split(input.Pictures, ","), UserId: input.UserId}
		AdID, err := domain.NewAdID()
		if err != nil {
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
