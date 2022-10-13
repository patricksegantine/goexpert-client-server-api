package entities

import "context"

// ExchangeRate holds exchange rate data
type ExchangeRate struct {
	ID         int64
	Name       string
	Bid        string
	CreateDate string
}

type ExchangeRateRepository interface {
	// Create create a new record in database
	Create(ctx context.Context, er ExchangeRate) error
}
