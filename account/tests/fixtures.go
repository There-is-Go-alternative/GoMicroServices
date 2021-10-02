package tests

import (
	"context"
	"fmt"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/account/usecase"
	"github.com/rs/zerolog/log"
)

var accountFixtures = []*domain.Account{
	{Firstname: "Damien", Lastname: "Bernard", Email: "damien.bernard@epitech.eu", Admin: true,
		Address: domain.Address{Country: "FR", State: "IDF", City: "Paris", Street: "Bleue", StreetNumber: 27}},
	{Firstname: "Naoufel", Lastname: "BERRADA", Email: "naoufel.berrada@epitech.eu", Admin: true},
	{Firstname: "Okil saber", Lastname: "LAKHDARI", Email: "okil-saber.lakhdari@epitech.eu", Admin: true},
	{Firstname: "Anton", Lastname: "CAZALET", Email: "anton.cazalet@epitech.eu", Admin: true},
	{Firstname: "Lina", Lastname: "KACI", Email: "lina1.kaci@epitech.eu", Admin: true},
	{Firstname: "Nicolas", Lastname: "SARKOZY", Email: "en.prison@epitech.eu", Admin: false},
	{Firstname: "Ugo", Lastname: "LEVI--CESCUTTI", Email: "le.gros.connard@epitech.eu", Admin: false},
	{Firstname: "Slohan", Lastname: "SAINTE-CROIX", Email: "il.est.où.mon.punch.coco@epitech.eu", Admin: false},
	{Firstname: "Matias", Lastname: "CAMPOS", Email: "l'handicapé@epitech.eu", Admin: false},
}

func DefaultFixtures(ctx context.Context, db usecase.Database) error {
	return CreateFixtures(ctx, accountFixtures, db)
}

func CreateFixtures(ctx context.Context, fixtures []*domain.Account, db usecase.Database) error {
	for _, a := range fixtures {
		id, err := domain.NewAccountID()
		if err != nil {
			return fmt.Errorf("could not create fixtures: %v", err)
		}
		a.ID = *id
	}
	log.Info().Msgf("loading fixtures ...")
	return db.Create(ctx, fixtures...)
}
