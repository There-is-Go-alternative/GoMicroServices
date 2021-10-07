package usecase

import (
	"context"
	"log"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/database"
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authUseCase struct {
	db database.Database
}

type LoginProto func(ctx context.Context, input LoginDTO) (*domain.Token, error)

func Login(collection *mongo.Collection) LoginProto {
	return func(ctx context.Context, input LoginDTO) (*domain.Token, error) {
		var auth domain.Auth
		auth, err := database.FindByEmail(ctx, input.Email, collection)
		if err != nil {
			log.Println("Email or Password are incorrect")
			return nil, err
		}

		hashed, err := domain.HashPassword(input.Password)
		if err != nil {
			log.Println("Email or Password are incorrect")
			return nil, err
		}
		err = domain.VerifyPassword(hashed, auth.Password)
		if err != nil {
			log.Println("Email or Password are incorrect")
			return nil, err
		}

		token, err := domain.CreateToken(auth.ID)
		if err != nil {
			log.Println("An error occured while creating token")
			return nil, err
		}

		return &domain.Token{
			Token: token,
		}, nil
	}
}
