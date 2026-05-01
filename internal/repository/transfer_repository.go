package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type transferRepository struct {
	pool *pgxpool.Pool
}

type account struct {
	id          int
	currency_id int
	balance     decimal.Decimal
	is_active   bool
}

func NewTransferRepository(pool *pgxpool.Pool) TransferRepository {
	return &transferRepository{pool: pool}
}
