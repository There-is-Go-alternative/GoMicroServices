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
		ad_list, err := u.DB.All(ctx)

		if err != nil {
			return nil, err
		}
		new_ad_list := make([]*domain.Ad, 0)

		for _, ad := range ad_list {
			if strings.Contains(strings.ToLower(ad.Title), content) {
				 new_ad_list = append(new_ad_list, ad)
			} else if strings.Contains(strings.ToLower(ad.Description), content) {
				new_ad_list = append(new_ad_list, ad)
		   }
		}
		return new_ad_list, nil
	}
}
