package database

import (
	"context"
	"fmt"
	"log"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type MongoDB struct {
	collection *mongo.Collection
	client     *mongo.Client
}

func NewConnection(dbName string, collectionName string, mongoURI string) (*MongoDB, error) {
	client, err := GetMongoDbConnection(mongoURI)
	if err != nil {
		return nil, err
	}
	collection, err := GetMongoDbCollection(client, dbName, collectionName)
	if err != nil {
		return nil, err
	}
	return &MongoDB{
		collection: collection,
		client:     client,
	}, nil
}

func GetMongoDbCollection(client *mongo.Client, dbName string, collectionName string) (*mongo.Collection, error) {
	collection := client.Database(dbName).Collection(collectionName)
	return collection, nil
}

func GetMongoDbConnection(mongoURI string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("Database mongoDB failed to connect: %v", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("Database mongoDB failed to ping: %v", err)
	}

	return client, nil
}

func (db *MongoDB) FindByEmail(ctx context.Context, email string) (domain.Auth, error) {
	var auth domain.Auth
	err := db.collection.FindOne(ctx, bson.M{"email": email}).Decode(&auth)
	if err != nil {
		log.Println(err)
		return domain.Auth{}, err
	}
	return auth, nil
}

func (db *MongoDB) FindByID(ctx context.Context, id string) (domain.Auth, error) {
	var auth domain.Auth
	err := db.collection.FindOne(ctx, bson.M{"id": id}).Decode(&auth)
	if err != nil {
		log.Println(err)
		return domain.Auth{}, err
	}
	return auth, nil
}
