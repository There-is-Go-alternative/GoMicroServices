package main

import (
	"log"
	"os"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/database"
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/transport/http"
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	client, err := database.GetMongoDbConnection()
	if err != nil {
		log.Fatal(err)
		return
	}
	db, err := database.GetMongoDbCollection(client, os.Getenv("MONGO_DB"), os.Getenv("MONGO_COLLECTION"))
	if err != nil {
		log.Fatal(err)
		return
	}
	authUseCase := usecase.NewUseCase(db)
	server := http.NewHttpServer(authUseCase)
	errc := make(chan error)
	go func() {
		errc <- server.Listen(os.Getenv("PORT"))
	}()
	err = <-errc
	if err != nil {
		log.Fatal(err)
	}
}
