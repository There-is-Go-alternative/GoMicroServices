package usecase

import (
	"context"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal"
)

type SearchAdCmd func(ctx context.Context, input SearchAdInput) ([]domain.Ad, error)

type SearchAdInput struct {
	Content string `json:"content" binding:"required"`
}

func (u UseCase) SearchAd() SearchAdCmd {
	return func(ctx context.Context, input SearchAdInput) ([]domain.Ad, error) {
		content := input.Content

		res, err := u.DB.Search(ctx, content)
		if err != nil {
			return nil, internal.NewInternalError(internal.InternalServerError, internal.InternalServerErrorMsg)
		}

		return res, nil
	}
}
