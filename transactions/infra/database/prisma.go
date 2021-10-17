//go:generate go run github.com/prisma/prisma-client-go generate
package database

import (
	"context"
	"time"

	"github.com/There-is-Go-alternative/GoMicroServices/transactions/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/transactions/infra/database/db"
)

type PrismaDB struct {
	Client *db.PrismaClient
}

func NewPrismaDB() *PrismaDB {
	return &PrismaDB{
		Client: db.NewClient(),
	}
}

func (p PrismaDB) Connect() error {
	return p.Client.Connect()
}

func (p PrismaDB) Disconnect() error {
	return p.Client.Disconnect()
}

func (p PrismaDB) All(ctx context.Context) ([]*domain.Transaction, error) {
	transactions, err := p.Client.Transaction.FindMany().Exec(ctx)

	if err != nil {
		return nil, err
	}

	formatted := make([]*domain.Transaction, len(transactions))

	for i, transaction := range transactions {
		formatted[i] = &domain.Transaction{
			Id:       transaction.ID,
			SellerId: transaction.SellerID,
			BuyerId:  transaction.BuyerID,
			Date:     transaction.Date,
			AdId:     transaction.AdID,
			Price:    transaction.Price,
		}
	}
	return formatted, nil
}

func (p PrismaDB) ById(ctx context.Context, id string) (*domain.Transaction, error) {
	transaction, err := p.Client.Transaction.FindUnique(
		db.Transaction.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &domain.Transaction{
		Id:       transaction.ID,
		SellerId: transaction.SellerID,
		BuyerId:  transaction.BuyerID,
		Date:     transaction.Date,
		AdId:     transaction.AdID,
		Price:    transaction.Price,
	}, nil
}

func (p PrismaDB) ByAdId(ctx context.Context, id string) (*domain.Transaction, error) {
	transaction, err := p.Client.Transaction.FindUnique(
		db.Transaction.AdID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &domain.Transaction{
		Id:       transaction.ID,
		SellerId: transaction.SellerID,
		BuyerId:  transaction.BuyerID,
		Date:     transaction.Date,
		AdId:     transaction.AdID,
		Price:    transaction.Price,
	}, nil
}

func (p PrismaDB) ByDateRange(ctx context.Context, start time.Time, end time.Time) ([]*domain.Transaction, error) {
	transactions, err := p.Client.Transaction.FindMany(
		db.Transaction.Date.BeforeEquals(end),
		db.Transaction.Date.AfterEquals(start),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	formatted := make([]*domain.Transaction, len(transactions))

	for i, transaction := range transactions {
		formatted[i] = &domain.Transaction{
			Id:       transaction.ID,
			SellerId: transaction.SellerID,
			BuyerId:  transaction.BuyerID,
			Date:     transaction.Date,
			AdId:     transaction.AdID,
			Price:    transaction.Price,
		}
	}
	return formatted, nil
}

func (p PrismaDB) BySellerId(ctx context.Context, id string) ([]*domain.Transaction, error) {
	transactions, err := p.Client.Transaction.FindMany(
		db.Transaction.SellerID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	formatted := make([]*domain.Transaction, len(transactions))

	for i, transaction := range transactions {
		formatted[i] = &domain.Transaction{
			Id:       transaction.ID,
			SellerId: transaction.SellerID,
			BuyerId:  transaction.BuyerID,
			Date:     transaction.Date,
			AdId:     transaction.AdID,
			Price:    transaction.Price,
		}
	}
	return formatted, nil
}

func (p PrismaDB) ByBuyerId(ctx context.Context, id string) ([]*domain.Transaction, error) {
	transactions, err := p.Client.Transaction.FindMany(
		db.Transaction.BuyerID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	formatted := make([]*domain.Transaction, len(transactions))

	for i, transaction := range transactions {
		formatted[i] = &domain.Transaction{
			Id:       transaction.ID,
			SellerId: transaction.SellerID,
			BuyerId:  transaction.BuyerID,
			Date:     transaction.Date,
			AdId:     transaction.AdID,
			Price:    transaction.Price,
		}
	}
	return formatted, nil
}

func (p PrismaDB) ByBuyerIdDateRange(ctx context.Context, id string, start time.Time, end time.Time) ([]*domain.Transaction, error) {
	transactions, err := p.Client.Transaction.FindMany(
		db.Transaction.BuyerID.Equals(id),
		db.Transaction.Date.BeforeEquals(end),
		db.Transaction.Date.AfterEquals(start),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	formatted := make([]*domain.Transaction, len(transactions))

	for i, transaction := range transactions {
		formatted[i] = &domain.Transaction{
			Id:       transaction.ID,
			SellerId: transaction.SellerID,
			BuyerId:  transaction.BuyerID,
			AdId:     transaction.AdID,
			Date:     transaction.Date,
			Price:    transaction.Price,
		}
	}
	return formatted, nil
}

func (p PrismaDB) BySellerIdDateRange(ctx context.Context, id string, start time.Time, end time.Time) ([]*domain.Transaction, error) {
	transactions, err := p.Client.Transaction.FindMany(
		db.Transaction.SellerID.Equals(id),
		db.Transaction.Date.BeforeEquals(end),
		db.Transaction.Date.AfterEquals(start),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	formatted := make([]*domain.Transaction, len(transactions))

	for i, transaction := range transactions {
		formatted[i] = &domain.Transaction{
			Id:       transaction.ID,
			SellerId: transaction.SellerID,
			BuyerId:  transaction.BuyerID,
			AdId:     transaction.AdID,
			Date:     transaction.Date,
			Price:    transaction.Price,
		}
	}
	return formatted, nil
}

func (p PrismaDB) Register(ctx context.Context, t *domain.Transaction) error {
	if _, err := p.Client.Transaction.CreateOne(
		db.Transaction.ID.Set(t.Id),
		db.Transaction.SellerID.Set(t.SellerId),
		db.Transaction.BuyerID.Set(t.BuyerId),
		db.Transaction.AdID.Set(t.AdId),
		db.Transaction.Price.Set(t.Price),
	).Exec(ctx); err != nil {
		return err
	}

	return nil
}
