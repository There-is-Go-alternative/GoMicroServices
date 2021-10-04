package database

import (
	"context"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	db "github.com/There-is-Go-alternative/GoMicroServices/account/infra/database/prisma"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/xerrors"
)

type PrismaDB struct {
	Client *db.PrismaClient
}

// go run github.com/prisma/prisma-client-go db push --schema infra/database/prisma/schema.prisma

//https://github.com/prisma/prisma-client-go/blob/main/docs/quickstart.md

func NewPrismaDB() *PrismaDB {
	return &PrismaDB{
		Client: db.NewClient(),
	}
}

func prismaAccountToDomain(pa *db.AccountModel) *domain.Account {
	isAdmin := func() bool {
		return pa.Admin == db.RoleADMIN
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
		Balance:   pa.Balance,
		CreatedAt: pa.CreatedAt,
		UpdatedAt: pa.UpdatedAt,
	}
}

func (p PrismaDB) Create(ctx context.Context, accounts ...*domain.Account) error {
	var errs []error
	for _, a := range accounts {
		if _, err := p.Client.Account.CreateOne(
			db.Account.ID.Set(a.ID.String()),
			db.Account.Email.Set(a.Email),
			db.Account.Firstname.Set(a.Firstname),
			db.Account.Lastname.Set(a.Lastname),
			db.Account.Balance.Set(a.Balance),
			// TODO:
			//db.Account.Address.Link(),
		).Exec(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	return xerrors.Concat(errs...)
}

func (p PrismaDB) ByID(ctx context.Context, ID domain.AccountID) (*domain.Account, error) {
	a, err := p.Client.Account.FindUnique(db.Account.ID.Equals(ID.String())).Exec(ctx)
	if err != nil {
		return nil, xerrors.ErrorWithCode{Code: xerrors.ResourceNotFound, Err: err}
	}
	return prismaAccountToDomain(a), err
}

func (p PrismaDB) SearchBy(ctx context.Context, searchFuncs ...db.AccountWhereParam) ([]*domain.Account, error) {
	accountModels, err := p.Client.Account.FindMany(searchFuncs...).Exec(ctx)
	if err != nil {
		return nil, xerrors.ErrorWithCode{Code: xerrors.ResourceNotFound, Err: err}
	}
	accounts := make([]*domain.Account, len(accountModels))
	for _, a := range accountModels {
		acc := a
		accounts = append(accounts, prismaAccountToDomain(&acc))
	}
	return accounts, err
}

// ByEmail Retrieve the info that match "Email".
func (p PrismaDB) ByEmail(ctx context.Context, email string) ([]*domain.Account, error) {
	return p.SearchBy(ctx, db.Account.Email.Equals(email))
}

// ByFirstname Retrieve the info that match "FirstName".
func (p PrismaDB) ByFirstname(ctx context.Context, firstname string) ([]*domain.Account, error) {
	return p.SearchBy(ctx, db.Account.Firstname.Equals(firstname))
}

// ByLastname Retrieve the info that match "Lastname".
func (p PrismaDB) ByLastname(ctx context.Context, lastname string) ([]*domain.Account, error) {
	return p.SearchBy(ctx, db.Account.Firstname.Equals(lastname))
}

// ByFullname Retrieve the info that match "Firstname" and "Lastname".
func (p PrismaDB) ByFullname(ctx context.Context, firstname, lastname string) ([]*domain.Account, error) {
	return p.SearchBy(ctx,
		db.Account.Firstname.Equals(firstname),
		db.Account.Lastname.Equals(lastname),
	)
}

// Update a list of domain.Account to the MemMapStorage
func (p PrismaDB) Update(ctx context.Context, accounts ...*domain.Account) error {
	// TODO:
	return nil
	//var errs []error
	//for _, a := range accounts {
	//	updated, err := p.Client.Account.FindUnique(
	//		db.Account.ID.Equals(a.ID.String()),
	//	).Update(
	//
	//		Comment.Post.Link(
	//			Post.ID.Equals(postID),
	//		),
	//	).Exec(ctx)
	//	if err != nil {
	//		errs = append(errs, err)
	//	}
	//}
	//return p.Save(accounts...)
}

// All return all domain.Account.
func (p *PrismaDB) All(ctx context.Context) ([]*domain.Account, error) {
	return p.SearchBy(ctx)
}

// Remove a domain.Account from the MemMapStorage
func (p *PrismaDB) Remove(ctx context.Context, accounts ...*domain.Account) error {
	var errs []error
	for _, a := range accounts {
		_, err := p.Client.Account.FindUnique(db.Account.ID.Equals(a.ID.String())).Delete().Exec(ctx)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return xerrors.Concat(errs...)
}
