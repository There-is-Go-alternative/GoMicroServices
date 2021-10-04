package main

import (
	"fmt"

	account "github.com/There-is-Go-alternative/GoMicroServices/account/domain"
)

func main() {
	// memory := database.MockDatabase{}
	userA := account.Account{ID: "A"}
	userB := account.Account{ID: "B"}
	userC := account.Account{ID: "C"}
	fmt.Printf("%+v\n", userA.ID)
	fmt.Printf("%+v\n", userB.ID)
	fmt.Printf("%+v\n", userC.ID)
}
