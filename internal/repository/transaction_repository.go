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

func (tr *transactionRepository) Create(ctx context.Context, t ...domain.Transaction) ([]int, error) {
	q := querier(ctx, tr.pool)

	var id []int

	for _, val := range t {
		err := q.QueryRow(ctx, `insert into 
	transactions(type, amount, account_id, status, status_message) 
	values ($1,$2,$3,$4,$5) returning id`, val.Type, val.Amount, val.AccountId, val.Status, val.StatusMessage).Scan(&id)

		if err != nil {
			return []int{}, fmt.Errorf("Error occured while creating new transfer")
		}
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

// create table if not exists transactions(
//     id BIGINT generated always as identity PRIMARY KEY,
//     type transaction_type not null,
//     amount numeric(12,2) not null check (amount >0),
//     account_id BIGINT not null REFERENCES accounts(id),
//     status transaction_status not null DEFAULT 'pending',
//     status_message text not null,
//     created_at TIMESTAMPtz not null DEFAULT now()
// );

func (tr *transactionRepository) GetByAccountId(ctx context.Context, id int) ([]domain.Transaction, error) {
	rows, err := tr.pool.Query(ctx, `select id, type, amount, account_id, status, status_message, created_at from transactions where account_id = $1`, id)

	if err != nil {
		return []domain.Transaction{}, fmt.Errorf("Error in get transactions by id: %w", err)
	}

	sl := make([]domain.Transaction, 0)

	for rows.Next() {
		var t domain.Transaction
		err := rows.Scan(&t.Id, &t.Type, &t.Amount, &t.AccountId, &t.Status, &t.StatusMessage, &t.CreatedAt)
		if err != nil {
			return []domain.Transaction{}, fmt.Errorf("Error in a loop get transactions by id: %w", err)
		}
		sl = append(sl, t)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return []domain.Transaction{}, fmt.Errorf("Error in a row get transactions by id: %w", err)
	}

	return sl, nil
}

func (tr *transactionRepository) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	rows, err := tr.pool.Query(ctx, `select id, type, amount, account_id, status, status_message, created_at from transactions`)

	if err != nil {
		return []domain.Transaction{}, fmt.Errorf("Error in get transactions by id: %w", err)
	}

	sl := make([]domain.Transaction, 0)

	for rows.Next() {
		var t domain.Transaction
		err := rows.Scan(&t.Id, &t.Type, &t.Amount, &t.AccountId, &t.Status, &t.StatusMessage, &t.CreatedAt)
		if err != nil {
			return []domain.Transaction{}, fmt.Errorf("Error in a loop get transactions by id: %w", err)
		}
		sl = append(sl, t)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return []domain.Transaction{}, fmt.Errorf("Error in a row get transactions by id: %w", err)
	}

	return sl, nil
}
