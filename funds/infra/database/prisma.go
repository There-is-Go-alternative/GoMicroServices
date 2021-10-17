//go:generate go run github.com/prisma/prisma-client-go generate
package database

import (
	"context"
	"time"

	"github.com/There-is-Go-alternative/GoMicroServices/funds/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/funds/infra/database/db"
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

func (p PrismaDB) ByUserId(ctx context.Context, id string) (*domain.Funds, error) {
	fund, err := p.Client.Funds.FindUnique(
		db.Funds.UserID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &domain.Funds{
		ID:          domain.FundsID(fund.ID),
		UserId:      fund.UserID,
		Balance:     fund.Balance,
		LastUpdated: fund.LastUpdated,
	}, nil
}

func (p PrismaDB) ById(ctx context.Context, id *domain.FundsID) (*domain.Funds, error) {
	fund, err := p.Client.Funds.FindUnique(
		db.Funds.ID.Equals(id.String()),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &domain.Funds{
		ID:          domain.FundsID(fund.ID),
		UserId:      fund.UserID,
		Balance:     fund.Balance,
		LastUpdated: fund.LastUpdated,
	}, nil
}

func (p PrismaDB) Update(ctx context.Context, id *domain.FundsID, new_balance float64) error {
	_, err := p.Client.Funds.FindUnique(
		db.Funds.ID.Equals(id.String()),
	).Update(
		db.Funds.Balance.Set(new_balance),
		db.Funds.LastUpdated.Set(time.Now()),
	).Exec(ctx)

	return err
}

func (p PrismaDB) UpdateByUser(ctx context.Context, id string, new_balance float64) error {
	_, err := p.Client.Funds.FindUnique(
		db.Funds.UserID.Equals(id),
	).Update(
		db.Funds.Balance.Set(new_balance),
		db.Funds.LastUpdated.Set(time.Now()),
	).Exec(ctx)

	return err
}

func (p PrismaDB) IncreaseByUser(ctx context.Context, id string, new_balance float64) error {
	_, err := p.Client.Funds.FindUnique(
		db.Funds.UserID.Equals(id),
	).Update(
		db.Funds.Balance.Increment(new_balance),
		db.Funds.LastUpdated.Set(time.Now()),
	).Exec(ctx)

	return err
}

func (p PrismaDB) Increase(ctx context.Context, id *domain.FundsID, new_balance float64) error {
	_, err := p.Client.Funds.FindUnique(
		db.Funds.ID.Equals(id.String()),
	).Update(
		db.Funds.Balance.Increment(new_balance),
		db.Funds.LastUpdated.Set(time.Now()),
	).Exec(ctx)

	return err
}

func (p PrismaDB) DecreaseByUser(ctx context.Context, id string, new_balance float64) error {
	_, err := p.Client.Funds.FindUnique(
		db.Funds.UserID.Equals(id),
	).Update(
		db.Funds.Balance.Decrement(new_balance),
		db.Funds.LastUpdated.Set(time.Now()),
	).Exec(ctx)

	return err
}

func (p PrismaDB) Decrease(ctx context.Context, id *domain.FundsID, new_balance float64) error {
	_, err := p.Client.Funds.FindUnique(
		db.Funds.ID.Equals(id.String()),
	).Update(
		db.Funds.Balance.Decrement(new_balance),
		db.Funds.LastUpdated.Set(time.Now()),
	).Exec(ctx)

	return err
}

func (p PrismaDB) Create(ctx context.Context, f domain.Funds) error {
	if _, err := p.Client.Funds.CreateOne(
		db.Funds.ID.Set(f.ID.String()),
		db.Funds.UserID.Set(f.UserId),
		db.Funds.Balance.Set(f.Balance),
		db.Funds.LastUpdated.Set(f.LastUpdated),
	).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (p PrismaDB) DeleteByUserId(ctx context.Context, id string) error {
	_, err := p.Client.Funds.FindUnique(
		db.Funds.UserID.Equals(id),
	).Delete().Exec(ctx)

	return err
}

func (p PrismaDB) DeleteById(ctx context.Context, id *domain.FundsID) error {
	_, err := p.Client.Funds.FindUnique(
		db.Funds.ID.Equals(id.String()),
	).Delete().Exec(ctx)

	return err
}

func (p PrismaDB) All(ctx context.Context) ([]*domain.Funds, error) {
	list, err := p.Client.Funds.FindMany().Exec(ctx)

	if err != nil {
		return nil, err
	}

	fundsList := make([]*domain.Funds, 0, len(list))

	for _, element := range list {
		fundsList = append(fundsList, &domain.Funds{
			ID:          domain.FundsID(element.ID),
			UserId:      element.UserID,
			Balance:     element.Balance,
			LastUpdated: element.LastUpdated,
		})
	}

	return fundsList, nil
}
