package repository

import (
	"context"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transferRepository struct {
	pool *pgxpool.Pool
}

func NewTransferRepository(pool *pgxpool.Pool) TransferRepository {
	return &transferRepository{pool: pool}
}

// CREATE TABLE IF NOT EXISTS transfers(
//     id BIGINT generated always as identity PRIMARY key,
//     transaction_id BIGINT REFERENCES transactions(id),
//     sender_account_id BIGINT not null REFERENCES accounts(id),
//     receiver_account_id BIGINT not null REFERENCES accounts(id),
//     currency_id BIGINT not null REFERENCES currencies(id),
//     amount numeric(12,2) not null check (amount > 0)
// );

func (tr *transferRepository) Create(ctx context.Context, t ...domain.Transfer) error {
	q := querier(ctx, tr.pool)

	for _, val := range t {
		_, err := q.Exec(ctx, `insert into transfers(transaction_id, sender_account_id, receiver_account_id, amount, currency_id) values($1, $2,$3,$4, $5)`, val.TransactionId, val.SenderAccountId, val.ReceiverAccountId, val.Amount, val.CurrencyId)

		if err != nil {
			return fmt.Errorf("Error in creating transfer: %w", err)
		}
	}
	return nil
}
