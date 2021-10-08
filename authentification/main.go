package main

import (
	"os"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/database"
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/transport/http"
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := database.NewConnection(
		os.Getenv("MONGO_DB"),
		os.Getenv("MONGO_COLLECTION"),
		os.Getenv("MONGO_URI"))
	if err != nil {
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
	}
}
