package awesomeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type AwesomeApiClient interface {
	GetExchangeRate(ctx context.Context, currency string) (*ExchangeRateResponse, error)
}

type AwesomeApiClientImpl struct{}

func NewAwesomeApiClient() AwesomeApiClient {
	return new(AwesomeApiClientImpl)
}

// GetExchangeRate gets the exchange rate by currency pair, e.g. USD-BRL or GBP-BRL
func (api *AwesomeApiClientImpl) GetExchangeRate(ctx context.Context, currency string) (*ExchangeRateResponse, error) {
	c := http.Client{Timeout: 300 * time.Millisecond}
	resp, err := c.Get(fmt.Sprintf("https://economia.awesomeapi.com.br/json/last/%v", currency))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &ErrAwesomeApi{resp.Status, resp.StatusCode}
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(strings.NewReader(string(content)))

	var exchangeRate map[string]ExchangeRateResponse
	for {
		if err := dec.Decode(&exchangeRate); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	key := strings.ToUpper(strings.Replace(currency, "-", "", -1))

	ret := exchangeRate[key]
	return &ret, nil
}
