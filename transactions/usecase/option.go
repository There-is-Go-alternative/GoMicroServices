package usecase

import (
	"context"
	"time"

	"github.com/There-is-Go-alternative/GoMicroServices/transactions/domain"
)

type DateRangeInput struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type GetByDateRangeCmd func(ctx context.Context, start, end time.Time) ([]*domain.Transaction, error)

func (u UseCase) GetByDateRange() GetByDateRangeCmd {
	return func(ctx context.Context, start, end time.Time) ([]*domain.Transaction, error) {
		return u.DB.ByDateRange(ctx, start, end)
	}
}

type GetBySellerIdDateRangeCmd func(ctx context.Context, id string, start, end time.Time) ([]*domain.Transaction, error)

func (u UseCase) GetBySellerIdDateRange() GetBySellerIdDateRangeCmd {
	return func(ctx context.Context, id string, start, end time.Time) ([]*domain.Transaction, error) {
		return u.DB.BySellerIdDateRange(ctx, id, start, end)
	}
}

type GetByBuyerIdDateRangeCmd func(ctx context.Context, id string, start, end time.Time) ([]*domain.Transaction, error)

func (u UseCase) GetByBuyerIdDateRange() GetByBuyerIdDateRangeCmd {
	return func(ctx context.Context, id string, start, end time.Time) ([]*domain.Transaction, error) {
		return u.DB.ByBuyerIdDateRange(ctx, id, start, end)
	}
}
