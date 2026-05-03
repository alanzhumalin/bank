package repository

import (
	"context"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{
		pool: pool,
	}
}

func (tr *transactionRepository) Create(ctx context.Context, t domain.Transaction) (int, error) {
	q := querier(ctx, tr.pool)

	var id int

	err := q.QueryRow(ctx, `insert into 
	transactions(type, amount, account_id, status, status_message) 
	values ($1,$2,$3,$4,$5) returning id`, t.Type, t.Amount, t.AccountId, t.Status, t.StatusMessage).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("Error occured while creating new transfer")
	}

	return id, nil
}

func (tr *transactionRepository) MarkTransaction(ctx context.Context, status string, status_message string, id int) error {
	q := querier(ctx, tr.pool)

	commandTag, err := q.Exec(ctx, `update transactions set status = $1, status_message = $2 where id = $3`, status, status_message, id)

	if err != nil {
		return fmt.Errorf("Error in mark complete the transaction: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrorTransactionNotFound
	}

	return nil

}
