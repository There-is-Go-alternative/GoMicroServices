package usecase

import (
	"context"
	"log"
	"os"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/database"
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginProto func(ctx context.Context, input LoginDTO) (*domain.Token, error)

func Login() LoginProto {
	return func(ctx context.Context, input LoginDTO) (*domain.Token, error) {
		collection, err := database.GetMongoDbCollection(os.Getenv("MONGO_DB"), os.Getenv("MONGO_COLLECTION"))
		if err != nil {
			return nil, err
		}

		var auth domain.Auth
		err = collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&auth)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		var hashed = domain.HashPassword(input.Password)
		err = domain.VerifyPassword(hashed, auth.Password)
		if err != nil {
			return nil, err
		}

		token, err := domain.CreateToken(auth.ID)
		if err != nil {
			return nil, err
		}

		return &domain.Token{
			Token: token,
		}, nil
	}
}
