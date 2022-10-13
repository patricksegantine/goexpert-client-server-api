package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/patricksegantine/goexpert-client-server-api/internal/dto"
	"github.com/patricksegantine/goexpert-client-server-api/internal/entities"
	"github.com/patricksegantine/goexpert-client-server-api/internal/infra/external_services/awesomeapi"
)

type CotacaoHandler struct {
	client awesomeapi.AwesomeApiClient
	repo   entities.ExchangeRateRepository
}

func NewCotacaoHandler(client awesomeapi.AwesomeApiClient, repo entities.ExchangeRateRepository) *CotacaoHandler {
	return &CotacaoHandler{client: client, repo: repo}
}

func (h *CotacaoHandler) GetExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	currency := r.URL.Query().Get("currency")
	if currency == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exchange, err := h.client.GetExchangeRate(ctx, currency)
	if err != nil {
		if err, ok := err.(awesomeapi.ErrAwesomeApi); ok {
			if err.StatusCode == 404 {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}

		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error getting exchange rate: %v\n", err)
		return
	}

	err = h.repo.Create(ctx, entities.ExchangeRate{Name: exchange.Name, Bid: exchange.Bid, CreateDate: exchange.CreateDate})
	if err != nil {
		log.Printf("Error persisting data: %v\n", err)
	}

	w.Header().Add("Contet-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&dto.ApiResponse{Name: exchange.Name, Bid: exchange.Bid})
}
