package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"io/ioutil"
	"net/http"
	"net/url"
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

type RegisterPayload struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

func (a *AuthHTTP) Register(email, password string, id domain.AccountID) error {
	payloadBuf := new(bytes.Buffer)
	data := url.Values{
		"id":       {id.String()},
		"email":    {email},
		"password": {password},
	}
	if err := json.NewEncoder(payloadBuf).Encode(data); err != nil {
		return fmt.Errorf("AuthHTTP register: %v", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/register", a.url), payloadBuf)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil
	}
	resp, err := a.Do(req, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("AuthHTTP register: Status differs from expected: %v (%v)", resp.StatusCode, http.StatusCreated)
	}
	return nil
}

func (a *AuthHTTP) Unregister(email, password string, id domain.AccountID) error {
	payloadBuf := new(bytes.Buffer)
	data := url.Values{
		"id":       {id.String()},
		"email":    {email},
		"password": {password},
	}
	if err := json.NewEncoder(payloadBuf).Encode(data); err != nil {
		return fmt.Errorf("AuthHTTP unregister: %v", err)
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/unregister", a.url), payloadBuf)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil
	}
	resp, err := a.Do(req, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("AuthHTTP register: Status differs from expected: %v (%v)", resp.StatusCode, http.StatusCreated)
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
