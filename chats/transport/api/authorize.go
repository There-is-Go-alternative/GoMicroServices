package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/chats/internal/xerrors"
	"github.com/gin-gonic/gin"
)

func SendRequest(token string) (string, error) {
	request, err := http.NewRequest("GET", "http://localhost:7500/account/test", nil)

	if err != nil {
		return "", err
	}
	//TODO Add "Bearer token"
	request.Header.Add("Authorization", fmt.Sprintf("%s", token))

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}
	return string([]byte(body)), nil
}

func Authorize(c *gin.Context) (*domain.Account, error) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		return nil, xerrors.AuthorizationError
	}

	account, err := SendRequest(token)

	if err != nil {
		return nil, err
	}

	var new_account_response struct {
		Data domain.Account `json:"data"`
	}
	err = json.Unmarshal([]byte(account), &new_account_response)

	if err != nil {
		return nil, err
	}
	return &new_account_response.Data, nil
}
