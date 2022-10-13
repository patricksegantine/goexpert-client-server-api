package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/patricksegantine/goexpert-client-server-api/internal/dto"
)

const (
	URL_COTACAO_API  string = "http://localhost:8080/cotacao/?currency=%s"
	CURRENCY_DEFAULT string = "USD-BRL"
)

func main() {
	defer fmt.Println("Finalizando servico...")

	fmt.Println("Iniciando servico...")

	moeda := CURRENCY_DEFAULT
	if len(os.Args) > 1 {
		moeda = os.Args[1]
	}

	cotacao, err := obterCotacao(moeda)
	if err != nil {
		panic(err)
	}

	err = salvarCotacao(moeda, cotacao)
	if err != nil {
		panic(err)
	}
}

func obterCotacao(moeda string) (*dto.ApiResponse, error) {
	c := http.Client{Timeout: 300 * time.Millisecond}
	resp, err := c.Get(fmt.Sprintf(URL_COTACAO_API, moeda))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cotacao dto.ApiResponse
	err = json.Unmarshal(content, &cotacao)
	if err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func salvarCotacao(moeda string, cotacao *dto.ApiResponse) error {
	file, err := os.Create("../../cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo: %v\n", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%s: {%s}", moeda, cotacao.Bid))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao gravar os dados no arquivo: %v\n", err)
	}
	fmt.Println("Arquivo criado com sucesso!")

	return nil
}
