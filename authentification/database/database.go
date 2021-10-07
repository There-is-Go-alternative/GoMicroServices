package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type Database interface {
	FindByEmail(context.Context, string) (domain.Auth, error)
}

type MongoDB struct {
	collection *mongo.Collection
}

func GetMongoDbCollection(client *mongo.Client, DbName string, CollectionName string) (*mongo.Collection, error) {
	collection := client.Database(DbName).Collection(CollectionName)
	return collection, nil
}

func GetMongoDbConnection() (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err != nil {
		return nil, fmt.Errorf("Database mongoDB failed to connect: %v", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("Database mongoDB failed to ping: %v", err)
	}

	return client, nil
}

func FindByEmail(ctx context.Context, email string, collection *mongo.Collection) (domain.Auth, error) {
	var auth domain.Auth
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&auth)
	if err != nil {
		log.Println(err)
		return domain.Auth{}, err
	}
	return auth, nil
}
