package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type AuthHTTP struct {
	client *http.Client
	url    string
	apiKey string
}

type authToken struct {
	Token string `json:"token"`
}

type credentialsPayload struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func NewAuthHTTP(url, apiKey string) *AuthHTTP {
	return &AuthHTTP{client: http.DefaultClient, url: url, apiKey: apiKey}
}

func (a AuthHTTP) Authorize(token string) (domain.AccountID, error) {
	token = strings.ReplaceAll(token, "Bearer ", "")
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(&authToken{Token: token}); err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/authorize", a.url), payloadBuf)
	if err != nil {
		return "", nil
	}

	rep := struct {
		UserID domain.AccountID `json:"user_id,required"`
	}{}

	resp, err := a.Do(req, &rep)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth: Status differs from expected: %v", http.StatusOK)
	}
	return rep.UserID, nil
}

func (a *AuthHTTP) Register(email, password string) error {
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(&credentialsPayload{Email: email, Password: password}); err != nil {
		return fmt.Errorf("AuthHTTP: %v", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/register", a.url), payloadBuf)
	if err != nil {
		return nil
	}
	rep := struct {
		domain.AccountID `json:"user_id,required"`
	}{}
	resp, err := a.Do(req, &rep)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth: Status differs from expected: %v", http.StatusOK)
	}
	return nil
}

func (a *AuthHTTP) Unregister(email string) error {
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(&credentialsPayload{Email: email}); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/unregister", a.url), payloadBuf)
	if err != nil {
		return nil
	}
	rep := struct {
		domain.AccountID `json:"user_id,required"`
	}{}
	resp, err := a.Do(req, &rep)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth: Status differs from expected: %v", http.StatusOK)
	}
	return nil
}

func (a *AuthHTTP) Do(r *http.Request, dst interface{}) (*http.Response, error) {
	response, err := a.client.Do(r)
	if err != nil {
		return nil, err
	}
	if dst == nil {
		return response, nil
	}

	defer func() { _ = response.Body.Close() }()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, dst); err != nil {
		return nil, err
	}
	return response, nil
}

func (a *AuthHTTP) GetAuthToken(ctx context.Context, userID string) (domain.Address, error) {
	uri := fmt.Sprintf("%s/auth/%s", a.url, userID)

	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return domain.Address{}, errors.Wrap(err, "unable to build the http request")
	}

	request = request.WithContext(ctx)
	request.Header.Set("api-key", a.apiKey)

	resp, err := a.client.Do(request)
	if err != nil {
		return domain.Address{}, errors.Wrap(err, "unable to handle the request")
	}
	switch resp.StatusCode {
	case http.StatusOK:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return domain.Address{}, errors.Wrap(err, "can't read body")
		}

		userAddress := domain.Address{}
		fmt.Println(string(data))
		err = json.Unmarshal(data, &userAddress)
		if err != nil {
			return domain.Address{}, errors.Wrap(err, "can't unmarshall json body")
		}
		fmt.Println(userAddress)
		return userAddress, nil
	default:
		return domain.Address{}, fmt.Errorf("GET user address API did not respond OK. HTTP code: %d", resp.StatusCode)
	}
}
