package awesomeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/patricksegantine/goexpert-client-server-api/internal/dto"
)

type ErrAwesomeApi struct {
	Status     string
	StatusCode int
}

func (e ErrAwesomeApi) Error() string {
	return e.Status
}

type AwesomeApiClient struct {
}

func NewAwesomeApiClient() *AwesomeApiClient {
	return new(AwesomeApiClient)
}

func (api *AwesomeApiClient) GetExchangeRate(moeda string) (*dto.CambioDto, error) {
	resp, err := http.Get(fmt.Sprintf("https://economia.awesomeapi.com.br/json/last/%v", moeda))
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

	var cambio map[string]dto.CambioDto
	for {
		if err := dec.Decode(&cambio); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	key := strings.ToUpper(strings.Replace(moeda, "-", "", -1))

	ret := cambio[key]
	return &ret, nil
}
