package usecase

import (
	"context"
	"fmt"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
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
			fmt.Println(err)
			return nil, err
		}
		return res, nil
	}
}
