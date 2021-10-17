package http

import (
	"encoding/json"
	"fmt"
	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"io/ioutil"
	"net/http"
)

type FundsHTTP struct {
	client *http.Client
	url    string
	apiKey string
}

func NewFundsHTTP(url, apiKey string) *FundsHTTP {
	return &FundsHTTP{client: http.DefaultClient, url: url, apiKey: apiKey}
}

func (a FundsHTTP) Create(ID domain.AccountID) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/%v", a.url, ID.String()), nil)
	if err != nil {
		return nil
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", a.apiKey))

	resp, err := a.Do(req, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("funds Create: Status differs from expected: %v", http.StatusOK)
	}
	return nil
}

func (a FundsHTTP) GetByID(ID domain.AccountID) (*float64, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/user/%v", a.url, ID.String()), nil)
	if err != nil {
		return nil, nil
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", a.apiKey))

	rep := domain.Balance{}
	resp, err := a.Do(req, &rep)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("funds Get: Status differs from expected: %v", http.StatusOK)
	}
	return &rep.Balance, nil
}

func (a FundsHTTP) GetAll() (map[domain.AccountID]*domain.Balance, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/%v", a.url), nil)
	if err != nil {
		return nil, nil
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", a.apiKey))

	var rep []*domain.Balance
	resp, err := a.Do(req, &rep)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("funds Get: Status differs from expected: %v", http.StatusOK)
	}
	m := make(map[domain.AccountID]*domain.Balance)
	for _, b := range rep {
		m[domain.AccountID(b.UserId)] = b
	}
	return m, nil
}

func (a FundsHTTP) Delete(ID domain.AccountID) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/%v", a.url, ID.String()), nil)
	if err != nil {
		return nil
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", a.apiKey))

	resp, err := a.Do(req, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth: Status differs from expected: %v", http.StatusOK)
	}
	return nil
}

func (a *FundsHTTP) Do(r *http.Request, dst interface{}) (*http.Response, error) {
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
