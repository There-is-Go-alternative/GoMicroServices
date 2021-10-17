package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/transactions/domain"
	"github.com/pkg/errors"
)

type FundsHTTP struct {
	client   *http.Client
	url      string
	apiToken string
}

func NewFundsHTTP(url, apiToken string) *FundsHTTP {
	return &FundsHTTP{client: http.DefaultClient, url: url, apiToken: apiToken}
}

type updateBalanceBody struct {
	By float64 `json:"by" binding:"required"`
}

func (f FundsHTTP) Increase(ctx context.Context, user_id string, by float64) error {
	uri := fmt.Sprintf("%s/api/v1/balance/increase/user/%s", f.url, user_id)

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(&updateBalanceBody{By: by}); err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, uri, payloadBuf)
	if err != nil {
		return errors.Wrap(err, "unable to build the http request")
	}

	request = request.WithContext(ctx)
	request.Header.Set("Authorize", fmt.Sprintf("Bearer %s", f.apiToken))

	resp, err := f.client.Do(request)

	if err != nil {
		return errors.Wrap(err, "unable to handle the request")
	}
	switch resp.StatusCode {
	case http.StatusOK:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "can't read body")
		}

		ad := &domain.Ad{}
		fmt.Println(string(data))
		err = json.Unmarshal(data, &ad)
		if err != nil {
			return errors.Wrap(err, "can't unmarshall json body")
		}
		return nil
	default:
		return fmt.Errorf("GET user address API did not respond OK. HTTP code: %d", resp.StatusCode)
	}
}

func (f FundsHTTP) Decrease(ctx context.Context, user_id string, by float64) error {
	uri := fmt.Sprintf("%s/api/v1/balance/decrease/user/%s", f.url, user_id)

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(&updateBalanceBody{By: by}); err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, uri, payloadBuf)
	if err != nil {
		return errors.Wrap(err, "unable to build the http request")
	}

	request = request.WithContext(ctx)
	request.Header.Set("Authorize", fmt.Sprintf("Bearer %s", f.apiToken))

	resp, err := f.client.Do(request)

	if err != nil {
		return errors.Wrap(err, "unable to handle the request")
	}
	switch resp.StatusCode {
	case http.StatusOK:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "can't read body")
		}

		ad := &domain.Ad{}
		fmt.Println(string(data))
		err = json.Unmarshal(data, &ad)
		if err != nil {
			return errors.Wrap(err, "can't unmarshall json body")
		}
		return nil
	default:
		return fmt.Errorf("GET user address API did not respond OK. HTTP code: %d", resp.StatusCode)
	}
}

func (f FundsHTTP) GetBalance(ctx context.Context, user_id string) (*float64, error) {
	uri := fmt.Sprintf("%s/api/v1/balance/user/%s", f.url, user_id)

	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build the http request")
	}

	request = request.WithContext(ctx)
	request.Header.Set("Authorize", fmt.Sprintf("Bearer %s", f.apiToken))

	resp, err := f.client.Do(request)

	if err != nil {
		return nil, errors.Wrap(err, "unable to handle the request")
	}
	switch resp.StatusCode {
	case http.StatusOK:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "can't read body")
		}

		funds := &domain.Funds{}
		fmt.Println(string(data))
		err = json.Unmarshal(data, &funds)
		if err != nil {
			return nil, errors.Wrap(err, "can't unmarshall json body")
		}
		return &funds.Balance, nil
	default:
		return nil, fmt.Errorf("GET user address API did not respond OK. HTTP code: %d", resp.StatusCode)
	}
}
