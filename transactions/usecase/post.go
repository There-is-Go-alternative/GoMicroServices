package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/There-is-Go-alternative/GoMicroServices/transactions/domain"
)

type getUser func(ctx context.Context, token string) (*domain.Account, error)

func getUserId(ctx context.Context, token string, GetUser getUser) (*string, error) {
	if token == "" {
		return nil, fmt.Errorf("no token provided")
	}
	hasPrefix := strings.HasPrefix(token, "Bearer ")

	if !hasPrefix {
		return nil, fmt.Errorf("wrong token format")
	}
	token = token[7:]
	userId, err := GetUser(ctx, token)

	if err != nil {
		return nil, err
	}
	return &userId.Id, nil
}

type RegiterCmd func(ctx context.Context, ad_id, token string) error

func (u UseCase) Register() RegiterCmd {
	return func(ctx context.Context, ad_id, token string) error {
		ad, err := u.Ads.GetAd(ctx, ad_id)

		if err != nil {
			return err
		}

		buyer_id, err := getUserId(ctx, token, u.Auth.GetUser)

		if err != nil {
			return err
		}

		if id, err := domain.NewTransactionID(); err != nil {
			return err
		} else {
			transaction := &domain.Transaction{
				Id:       *id,
				SellerId: ad.UserId,
				BuyerId:  *buyer_id,
				AdId:     ad_id,
				Date:     time.Now(),
				Price:    float64(ad.Price),
			}

			balance, err := u.Funds.GetBalance(ctx, *buyer_id)

			if err != nil {
				return err
			}

			if *balance < transaction.Price {
				return fmt.Errorf("user balance isn't filled enough")
			}

			err = u.Ads.Buy(ctx, ad_id, token)

			if err != nil {
				return err
			}

			if err = u.Funds.Decrease(ctx, transaction.BuyerId, transaction.Price); err != nil {
				return err
			}

			if err = u.Funds.Increase(ctx, transaction.SellerId, transaction.Price); err != nil {
				return err
			}

			return u.DB.Register(ctx, transaction)
		}
	}
}
