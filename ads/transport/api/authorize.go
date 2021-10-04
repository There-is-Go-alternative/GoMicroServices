package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/gin-gonic/gin"
)

func SendRequest(token string) string {
	request, err := http.NewRequest("GET", "http://localhost:7500/account/test", nil)

	if err != nil {
		return ""
	}
	request.Header.Add("Authorization", fmt.Sprintf("%s", token))

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

func Authorize(c *gin.Context) *domain.Account {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		return nil
	}
	//Todo call the account service (fix error)
	account := SendRequest("0647a8a4-a743-46de-bbc1-a87f8f1e0e43")

	if account == "" {
		//TODO return error
		return nil
	}
	var new_account_response struct {Data domain.Account `json:"data"`}

	err := json.Unmarshal([]byte(account), &new_account_response)

	if err != nil {
		//TODO Error
		return nil
	}
	return &new_account_response.Data
}
