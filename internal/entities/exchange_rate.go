package entities

type ExchangeRate struct {
	Name string
	Bid  string
}

func NewExchangeRate(name, bid string) *ExchangeRate {
	return &ExchangeRate{
		Name: name,
		Bid:  bid,
	}
}

type ExchangeRateRepository interface {
	Save(ExchangeRate) error
}
