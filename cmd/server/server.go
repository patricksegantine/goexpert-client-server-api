package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/patricksegantine/goexpert-client-server-api/cmd/server/handlers"
	"github.com/patricksegantine/goexpert-client-server-api/internal/entities"
	"github.com/patricksegantine/goexpert-client-server-api/internal/infra/database"
	"github.com/patricksegantine/goexpert-client-server-api/internal/infra/external_services/awesomeapi"
	"github.com/patricksegantine/goexpert-client-server-api/internal/infra/repositories"
)

var (
	client awesomeapi.AwesomeApiClient
	repo   entities.ExchangeRateRepository
)

/*
	Usando o package "context", o server.go deverá:
	1-registrar no banco de dados SQLite cada cotação recebida,
	sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e
	o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.
*/

func main() {
	db, err := sql.Open("sqlite3", "../../cotacao.db")
	if err != nil {
		log.Fatalf("Error initializing database connection: %v", err)
	}
	defer db.Close()

	err = database.SetupDatabase(db)
	if err != nil {
		log.Fatalf("Error checking database structure: %v", err)
	}

	repo = repositories.NewExchangeRateRepository(db)
	client = awesomeapi.NewAwesomeApiClient()
	handler := handlers.NewCotacaoHandler(client, repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is running. Request processed at %v", time.Now())
	})
	mux.HandleFunc("/cotacao/", handler.GetExchangeRateHandler)

	log.Println("Server running and listening at http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
