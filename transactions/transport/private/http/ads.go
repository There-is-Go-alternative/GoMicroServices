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

type AdsHTTP struct {
	client   *http.Client
	url      string
	apiToken string
}

func NewAdsHTTP(url, apiToken string) *AdsHTTP {
	return &AdsHTTP{client: http.DefaultClient, url: url, apiToken: apiToken}
}

func (a AdsHTTP) GetAd(ctx context.Context, token string) (*domain.Ad, error) {
	uri := fmt.Sprintf("%s/ads/", a.url)

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

		ad := &domain.Ad{}
		fmt.Println(string(data))
		err = json.Unmarshal(data, &ad)
		if err != nil {
			return nil, errors.Wrap(err, "can't unmarshall json body")
		}
		return ad, nil
	default:
		return nil, fmt.Errorf("GET user address API did not respond OK. HTTP code: %d", resp.StatusCode)
	}
}

func (a AdsHTTP) Buy(ctx context.Context, id string, token string) error {
	uri := fmt.Sprintf("%s/ads/buy/%s", a.url, id)

	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return errors.Wrap(err, "unable to build the http request")
	}

	request = request.WithContext(ctx)
	request.Header.Set("Authorize", fmt.Sprintf("Bearer %s", token))

	resp, err := a.client.Do(request)

	if err != nil {
		return errors.Wrap(err, "unable to handle the request")
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		return fmt.Errorf("GET user address API did not respond OK. HTTP code: %d", resp.StatusCode)
	}
}
