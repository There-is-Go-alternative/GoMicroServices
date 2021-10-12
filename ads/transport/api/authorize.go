package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal/xerrors"
	"github.com/gin-gonic/gin"
)

type DataElement struct{
    UserID string `json:"user_id"`
}

type Response struct {
    Data DataElement `json:"data"`
}

func SendRequest(token string) (string, error) {
	form := url.Values{}
	form.Set("token", token)
	request, err := http.NewRequest("POST", "http://localhost:8080/authorize", strings.NewReader(form.Encode()))

	if err != nil {
		return "", err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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

func Authorize(c *gin.Context) (string, error) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		return "", xerrors.AuthorizationError
	}

	account, err := SendRequest(token)

	if err != nil {
		return "", err
	}

	var new_account_response Response
	err = json.Unmarshal([]byte(account), &new_account_response)
	if err != nil {
		return "", err
	}
	
	if new_account_response.Data.UserID == "" {
		return "", xerrors.AuthorizationError
	}
	return new_account_response.Data.UserID, nil
}
