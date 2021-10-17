package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/transactions/domain"
	"github.com/pkg/errors"

	"io/ioutil"
)

type AccountHTTP struct {
	client *http.Client
	url    string
}

func NewAccountHTTP(url string) *AccountHTTP {
	return &AccountHTTP{client: http.DefaultClient, url: url}
}

func (a AccountHTTP) GetUser(ctx context.Context, token string) (*domain.Account, error) {
	uri := fmt.Sprintf("%s/", a.url)

	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build the http request")
	}

	request = request.WithContext(ctx)
	request.Header.Set("Authorize", fmt.Sprintf("Bearer %s", token))

	resp, err := a.client.Do(request)

	if err != nil {
		return nil, errors.Wrap(err, "unable to handle the request")
	}
	switch resp.StatusCode {
	case http.StatusOK:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "can't read body")
		}

		account := &domain.Account{}
		fmt.Println(string(data))
		err = json.Unmarshal(data, &account)
		if err != nil {
			return nil, errors.Wrap(err, "can't unmarshall json body")
		}
		return account, nil
	default:
		return nil, fmt.Errorf("GET user address API did not respond OK. HTTP code: %d", resp.StatusCode)
	}
}

func (a AccountHTTP) IsAdmin(id string) (bool, error) {
	return false, nil
}
