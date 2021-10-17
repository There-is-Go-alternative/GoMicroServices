package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal/xerrors"
	"github.com/imdario/mergo"
)

type UpdateAdCmd func(ctx context.Context, input UpdateAdInput) (*domain.Ad, error)

type UpdateAdInput struct {
	domain.Ad
}

func (u UseCase) UpdateAd() UpdateAdCmd {
	return func(ctx context.Context, input UpdateAdInput) (*domain.Ad, error) {
		ad, err := u.DB.ByID(ctx, input.ID)

		if err != nil {
			return nil, err
		}

		err = mergo.Merge(ad, &input.Ad, mergo.WithOverride)
		if err != nil {
			return nil, xerrors.ErrorWithCode{Code: xerrors.CodeInvalidData, Err: err}
		}

		err = u.DB.Update(ctx, ad)

		if err != nil {
			return nil, xerrors.ErrorWithCode{Code: xerrors.CodeInvalidData, Err: err}
		}

		return ad, err
	}
}
