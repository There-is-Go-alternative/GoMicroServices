package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal/xerrors"
)

type UpdateAdCmd func(ctx context.Context, input UpdateAdInput) (*domain.Ad, error)

type UpdateAdInput struct {
	ID domain.AdID `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Price uint `json:"price,omitempty"`
	Pictures string `json:"pictures,omitempty"`
}

func (u UseCase) UpdateAd() UpdateAdCmd {
	return func(ctx context.Context, input UpdateAdInput) (*domain.Ad, error) {
		ad := &domain.Ad{ID: domain.AdID(input.ID), Title: input.Title, Description: input.Description, Price: input.Price, Pictures: strings.Split(input.Pictures, ",")}

		if !ad.Validate() {
			return nil, xerrors.ErrorWithCode{
				Code: xerrors.CodeInvalidData, Err: fmt.Errorf("invalid user ad data: %v", ad),
			}
		}

		acc, _ := u.DB.ByID(ctx, ad.ID)

		if acc == nil {
			return nil, xerrors.ErrorWithCode{
				Code: xerrors.CodeInvalidData, Err: fmt.Errorf("account doesn't exists: %v", ad),
			}
		}
		err := u.DB.Update(ctx, ad)
		return ad, err
	}
}
