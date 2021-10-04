package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

//TODO: Remplacer valeur de retour par alias domain.Account
func SendRequest(token string) string {
	request, err := http.NewRequest("GET", "http://localhost:7500/", nil)

	if err != nil {
		return ""
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("TODO: Erreur")
		return ""
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Printf("TODO: Erreur")
		return ""
	}
	return string([]byte(body))
}

func Authorize(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		return ""
	}
	//Todo call the account service (fix error)
	account := SendRequest("TEST_TOKEN_ADMIN_INTRA_SERVICE")

	if account == "" {
		//TODO return error
		return ""
	}
	return account
}
