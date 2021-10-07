package main

import (
	"log"
	"os"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/database"
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/transport/http"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	client, err := database.GetMongoDbConnection()
	if err != nil {
		log.Fatal(err)
		return
	}
	collection, err := database.GetMongoDbCollection(client, os.Getenv("MONGO_DB"), os.Getenv("MONGO_COLLECTION"))
	if err != nil {
		log.Fatal(err)
		return
	}

	server := http.NewHttpServer(collection)
	errc := make(chan error)
	go func() {
		errc <- server.Listen(os.Getenv("PORT"))
	}()
	err = <-errc
	if err != nil {
		log.Fatal(err)
	}
}
