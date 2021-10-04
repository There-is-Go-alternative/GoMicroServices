package database

import (
	"context"
	"fmt"

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
	return nil, fmt.Errorf("Unimplemented")
}

func (p PrismaDB) ById(ctx context.Context, id *domain.FundsID) (*domain.Funds, error) {
	return nil, fmt.Errorf("Unimplemented")
}

func (p PrismaDB) Increase(ctx context.Context, by int) error {
	return fmt.Errorf("Unimplemented")
}

func (p PrismaDB) Decrease(ctx context.Context, by int) error {
	return fmt.Errorf("Unimplemented")
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
	return fmt.Errorf("Unimplemented")
}

func (p PrismaDB) DeleteById(ctx context.Context, id *domain.FundsID) error {
	return fmt.Errorf("Unimplemented")
}

func (p PrismaDB) All(ctx context.Context) ([]*domain.Funds, error) {
	list := p.Client.Funds.FindMany()

	fmt.Println(list)

	return make([]*domain.Funds, 0), nil
}
