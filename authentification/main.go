package main

import (
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/transport/http"
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	http.NewHttpServer(usecase.Login())
}
