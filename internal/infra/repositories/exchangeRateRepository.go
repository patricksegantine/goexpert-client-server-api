package repositories

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/patricksegantine/goexpert-client-server-api/internal/entities"
)

var (
	queries = map[string]string{
		"Create": `INSERT INTO exchange_rate(Name,Bid,CreateDate)
					VALUES (?,?,?);`,
	}
)

type ExchangeRateRepository struct {
	_  sync.Mutex
	db *sql.DB
}

func NewExchangeRateRepository(db *sql.DB) entities.ExchangeRateRepository {
	return &ExchangeRateRepository{db: db}
}

func (r *ExchangeRateRepository) Create(ctx context.Context, er entities.ExchangeRate) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	stmt, err := r.db.PrepareContext(ctx, queries["Create"])
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, er.Name, er.Bid, er.CreateDate)
	if err != nil {
		return err
	}
	return nil
}
