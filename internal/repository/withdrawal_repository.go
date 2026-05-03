package repository

import (
	"context"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type withdrawalRepository struct {
	pool *pgxpool.Pool
}

func NewWithdrawalRepository(pool *pgxpool.Pool) WithdrawalRepository {
	return &withdrawalRepository{
		pool: pool,
	}
}

func (wr *withdrawalRepository) Create(ctx context.Context, w domain.Withdrawal) error {
	q := querier(ctx, wr.pool)
	_, err := q.Exec(ctx, `insert into withdrawals(transaction_id, account_id, amount, source) 
	values($1, $2, $3, $4)`, w.TransactionId, w.AccountId, w.Amount, w.Source)

	if err != nil {
		return fmt.Errorf("Error in create withdrawal: %w", err)
	}
	return nil
}
