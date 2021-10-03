package usecase

import (
	"context"
	"strings"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/domain"
)

type SearchAdCmd func(ctx context.Context, input SearchAdInput) ([]*domain.Ad, error)

type SearchAdInput struct {
	Content string `json:"content" binding:"required"`
}

func (u UseCase) SearchAd() SearchAdCmd {
	return func(ctx context.Context, input SearchAdInput) ([]*domain.Ad, error) {
		content := strings.ToLower(input.Content)
		adList, err := u.DB.All(ctx)

		if err != nil {
			return nil, err
		}
		var newAdList []*domain.Ad

		for _, ad := range adList {
			if strings.Contains(strings.ToLower(ad.Title), content) {
				 newAdList = append(newAdList, ad)
			} else if strings.Contains(strings.ToLower(ad.Description), content) {
				newAdList = append(newAdList, ad)
		   }
		}
		return newAdList, nil
	}
}
