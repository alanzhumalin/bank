package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type idempotencyRepository struct {
	pool *pgxpool.Pool
}

func NewIdempotencyStore(pool *pgxpool.Pool) IdempotencyRepository {
	return &idempotencyRepository{
		pool: pool,
	}
}

func (i *idempotencyRepository) Exists(ctx context.Context, key string) (bool, error) {
	q := querier(ctx, i.pool)
	var ok bool
	err := q.QueryRow(ctx, `select exists(select 1 from idempotency_keys where idempotency_key = $1)`, key).Scan(&ok)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, fmt.Errorf("get idempotency key: %w", err)
	}

	if ok {
		return true, nil
	}

	return false, nil

}

func (i *idempotencyRepository) Start(ctx context.Context, idempotency domain.Idempotency) error {
	q := querier(ctx, i.pool)

	commandTag, err := q.Exec(ctx, `insert into idempotency_keys(user_id, idempotency_key, operation) values($1, $2, $3)`, idempotency.UserId, idempotency.IdempotencyKey, idempotency.Operation)

	if err != nil {
		return fmt.Errorf("Error in start of creating idempotency key: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrorNoRowsInserted
	}

	return nil
}

func (i *idempotencyRepository) Complete(ctx context.Context, idempotency domain.Idempotency) error {
	q := querier(ctx, i.pool)

	commandTag, err := q.Exec(ctx, `update idempotency_keys set transaction_id = $1, status = $2, response = $3,updated_at = $4 where idempotency_key = $5`,
		idempotency.TransactionId, idempotency.Status, idempotency.Response, idempotency.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("Error in completing of idempotency key: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrorNoRowsChanged
	}

	return nil
}
