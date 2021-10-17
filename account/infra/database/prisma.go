//go:generate go run github.com/prisma/prisma-client-go generate
package database

import (
	"context"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/account/infra/database/prisma"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/xerrors"
)

type DB struct {
	Client *prismaDB.PrismaClient
}

// go run github.com/prisma/prisma-client-go db push --schema infra/database/prisma/schema.prisma

//https://github.com/prisma/prisma-client-go/blob/main/docs/quickstart.md

func NewPrismaDB() (*DB, error) {
	db := prismaDB.NewClient()

	if err := db.Connect(); err != nil {
		return nil, err
	}

	return &DB{
		Client: db,
	}, nil
}

func prismaAccountToDomain(pa *prismaDB.AccountModel) *domain.Account {
	isAdmin := func() bool {
		return pa.Admin == prismaDB.RoleADMIN
	}
	getAddress := func() domain.Address {
		// TODO
		//if addr, ok := pa.Address(); !ok || addr == nil {
		//	return domain.Address{}
		//}
		return domain.Address{}
	}

	return &domain.Account{
		ID:        domain.AccountID(pa.ID),
		Email:     pa.Email,
		Firstname: pa.Firstname,
		Lastname:  pa.Lastname,
		Admin:     isAdmin(),
		Address:   getAddress(),
		CreatedAt: pa.CreatedAt,
		UpdatedAt: pa.UpdatedAt,
	}
}

func (p DB) Create(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	acc, err := p.Client.Account.CreateOne(
		prismaDB.Account.ID.Set(account.ID.String()),
		prismaDB.Account.Email.Set(account.Email),
		prismaDB.Account.Firstname.Set(account.Firstname),
		prismaDB.Account.Lastname.Set(account.Lastname),
		// TODO:
		//db.Account.Address.Link(),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return prismaAccountToDomain(acc), nil
}

func (p DB) ByID(ctx context.Context, ID domain.AccountID) (*domain.Account, error) {
	a, err := p.Client.Account.FindUnique(prismaDB.Account.ID.Equals(ID.String())).Exec(ctx)
	if err != nil {
		return nil, xerrors.ErrorWithCode{Code: xerrors.ResourceNotFound, Err: err}
	}
	return prismaAccountToDomain(a), err
}

func (p DB) SearchBy(ctx context.Context, searchFuncs ...prismaDB.AccountWhereParam) ([]*domain.Account, error) {
	accountModels, err := p.Client.Account.FindMany(searchFuncs...).Exec(ctx)
	if err != nil {
		return nil, xerrors.ErrorWithCode{Code: xerrors.ResourceNotFound, Err: err}
	}
	accounts := make([]*domain.Account, 0, len(accountModels))
	for _, a := range accountModels {
		acc := a
		accounts = append(accounts, prismaAccountToDomain(&acc))
	}
	return accounts, err
}

// ByEmail Retrieve the info that match "Email".
func (p DB) ByEmail(ctx context.Context, email string) ([]*domain.Account, error) {
	return p.SearchBy(ctx, prismaDB.Account.Email.Equals(email))
}

// ByFirstname Retrieve the info that match "FirstName".
func (p DB) ByFirstname(ctx context.Context, firstname string) ([]*domain.Account, error) {
	return p.SearchBy(ctx, prismaDB.Account.Firstname.Equals(firstname))
}

// ByLastname Retrieve the info that match "Lastname".
func (p DB) ByLastname(ctx context.Context, lastname string) ([]*domain.Account, error) {
	return p.SearchBy(ctx, prismaDB.Account.Firstname.Equals(lastname))
}

// ByFullname Retrieve the info that match "Firstname" and "Lastname".
func (p DB) ByFullname(ctx context.Context, firstname, lastname string) ([]*domain.Account, error) {
	return p.SearchBy(ctx,
		prismaDB.Account.Firstname.Equals(firstname),
		prismaDB.Account.Lastname.Equals(lastname),
	)
}

// Update a list of domain.Account to the MemMapStorage
func (p DB) Update(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	updated, err := p.Client.Account.FindUnique(
		prismaDB.Account.ID.Equals(account.ID.String()),
	).Update(
		prismaDB.Account.Email.Set(account.Email),
		prismaDB.Account.Firstname.Set(account.Firstname),
		prismaDB.Account.Lastname.Set(account.Lastname),
<<<<<<< HEAD
=======
		prismaDB.Account.Balance.Set(account.Balance),
>>>>>>> e6139e8 (add(account): Adding login and register management and cleaning errors.)
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return prismaAccountToDomain(updated), err
}

// All return all domain.Account.
func (p *DB) All(ctx context.Context) ([]*domain.Account, error) {
	return p.SearchBy(ctx)
}

// Remove a domain.Account from the MemMapStorage
func (p *DB) Remove(ctx context.Context, id domain.AccountID) (*domain.Account, error) {
	acc, err := p.Client.Account.FindUnique(prismaDB.Account.ID.Equals(id.String())).Delete().Exec(ctx)
	if err != nil {
		return nil, err
	}
	return prismaAccountToDomain(acc), nil
}
