package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/patricksegantine/goexpert-client-server-api/internal/infra/external_services/awesomeapi"
)

var (
	client *awesomeapi.AwesomeApiClient
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is running. Request processed at %v", time.Now())
	})
	mux.HandleFunc("/cotacao/", obterCotacaoHandler)

	client = awesomeapi.NewAwesomeApiClient()

	log.Println("Server running and listening at http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func obterCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	moeda := r.URL.Query().Get("moeda")
	if moeda == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cotacao, err := client.GetExchangeRate(moeda)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stdout, "Erro ao obter a cotacao: %v", err)
		return
	}

	w.Header().Add("Contet-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)
}
