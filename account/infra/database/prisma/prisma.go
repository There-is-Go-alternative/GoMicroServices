package prisma

import (
	"context"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/account/infra/database/prisma/prisma"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/xerrors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type DB struct {
	Client *prismaDB.PrismaClient
}

// go run github.com/prisma/prisma-client-go db push --schema infra/database/prisma/schema.prisma
// go run github.com/prisma/prisma-client-go generate --schema infra/database/prisma/schema.prisma

//https://github.com/prisma/prisma-client-go/blob/main/docs/quickstart.md

func NewPrismaDB() (*DB, error) {
	db := &DB{
		Client: prismaDB.NewClient(),
	}
	if err := db.Client.Prisma.Connect(); err != nil {
		return nil, err
	}

	return db, nil
}

func prismaAddressToDomain(pa *prismaDB.AddressModel) domain.Address {
	addr := domain.Address{
		Country:      pa.Country,
		State:        pa.State,
		City:         pa.City,
		Street:       pa.Street,
		StreetNumber: pa.StreetNumber,
	}
	if comp, ok := pa.Complementary(); ok {
		addr.Complementary = comp
	}
	return addr
}

func prismaAccountToDomain(pa *prismaDB.AccountModel) *domain.Account {
	isAdmin := func() bool {
		return pa.Admin == prismaDB.RoleADMIN
	}
	getAddress := func() domain.Address {
		// TODO
		if addr, ok := pa.Address(); ok && addr != nil {
			log.Warn().Msgf("fetching Address for account {%v}: %v", pa.Email, addr)
			return prismaAddressToDomain(addr)
		}
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

func (p DB) Create(ctx context.Context, accounts ...*domain.Account) error {
	var errs []error
	for _, a := range accounts {
		if _, err := p.Client.Account.CreateOne(
			prismaDB.Account.ID.Set(a.ID.String()),
			prismaDB.Account.Email.Set(a.Email),
			prismaDB.Account.Firstname.Set(a.Firstname),
			prismaDB.Account.Lastname.Set(a.Lastname),
			prismaDB.Account.Balance.Set(a.Balance),
		).Exec(ctx); err != nil {
			errs = append(errs, errors.Wrapf(err, "In PrismaDB.Create"))
		}

		if a.Address.Validate() != nil {
			log.Warn().Msgf("Address for account {%v} is not valid, skipping ...", a)
			continue
		}

		if _, err := p.Client.Address.CreateOne(
			prismaDB.Address.Account.Link(
				prismaDB.Account.ID.Equals(a.ID.String()),
			),
			prismaDB.Address.Country.Set(a.Address.Country),
			prismaDB.Address.State.Set(a.Address.State),
			prismaDB.Address.City.Set(a.Address.City),
			prismaDB.Address.Street.Set(a.Address.Street),
			prismaDB.Address.StreetNumber.Set(a.Address.StreetNumber),
		).Exec(ctx); err != nil {
			errs = append(errs, errors.Wrapf(err, "In PrismaDB.Create"))
		}
	}
	return xerrors.Concat(errs...)
}

func (p DB) ByID(ctx context.Context, ID domain.AccountID) (*domain.Account, error) {
	a, err := p.Client.Account.FindUnique(prismaDB.Account.ID.Equals(ID.String())).Exec(ctx)
	if err != nil {
		return nil, xerrors.ErrorWithCode{Code: xerrors.ResourceNotFound, Err: err}
	}
	return prismaAccountToDomain(a), err
}

func (p DB) SearchBy(ctx context.Context, searchFuncs ...prismaDB.AccountWhereParam) ([]*domain.Account, error) {
	accountModels, err := p.Client.Account.FindMany(
		searchFuncs...,
	).With(
	// TODO: fetch Address, pb in schema ?
	//db.Account.Address.Fetch()
	).Exec(ctx)
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
	return p.SearchBy(ctx,
		prismaDB.Account.Email.Equals(email),
	)
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
func (p DB) Update(ctx context.Context, accounts ...*domain.Account) error {
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
func (p *DB) All(ctx context.Context) ([]*domain.Account, error) {
	//return p.SearchBy(ctx)
	return p.SearchBy(ctx)
}

// Remove a domain.Account.
func (p *DB) Remove(ctx context.Context, accounts ...*domain.Account) error {
	var errs []error
	for _, a := range accounts {
		_, err := p.Client.Account.FindUnique(prismaDB.Account.ID.Equals(a.ID.String())).Delete().Exec(ctx)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return xerrors.Concat(errs...)
}
