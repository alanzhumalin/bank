package repository

import (
	"context"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type depositRepository struct {
	pool *pgxpool.Pool
}

func NewDepositRepository(pool *pgxpool.Pool) DepositRepository {
	return &depositRepository{
		pool: pool,
	}
}

func (dr *depositRepository) Create(ctx context.Context, d domain.Deposit) error {
	q := querier(ctx, dr.pool)

	_, err := q.Exec(ctx, `insert into deposits(transaction_id, account_id, amount, source) values($1,$2,$3,$4)`, d.TransactionId, d.AccountId, d.Amount, d.Source)

	if err != nil {
		return fmt.Errorf("Error in creating deposits: %w", err)
	}

	return nil
}
