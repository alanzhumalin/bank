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

func NewIdempotencyRepo(pool *pgxpool.Pool) IdempotencyRepository {
	return &idempotencyRepository{
		pool: pool,
	}
}

func (i *idempotencyRepository) GetByKey(ctx context.Context, key string, userId int) (domain.Idempotency, error) {
	q := querier(ctx, i.pool)
	var idempotency domain.Idempotency
	err := q.QueryRow(ctx, `select id, transaction_id, user_id, idempotency_key, operation, status, response, updated_at, created_at from idempotency_keys where user_id = $1 and idempotency_key = $2`, userId, key).Scan(
		&idempotency.Id, &idempotency.TransactionId, &idempotency.UserId, &idempotency.IdempotencyKey, &idempotency.Operation, &idempotency.Status, &idempotency.Response, &idempotency.UpdatedAt, &idempotency.CreatedAt,
	)

	if err == nil {
		return idempotency, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Idempotency{}, domain.ErrorIdempotencyKeyNotFound
	}

	return domain.Idempotency{}, fmt.Errorf("Error in getting idempotency key: %w", err)

}

func (i *idempotencyRepository) Start(ctx context.Context, idempotency domain.Idempotency) error {
	q := querier(ctx, i.pool)

	commandTag, err := q.Exec(ctx, `insert into idempotency_keys(user_id, idempotency_key, operation) values($1, $2, $3) on conflict (user_id, idempotency_key) do nothing`, idempotency.UserId, idempotency.IdempotencyKey, idempotency.Operation)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrorIdempotencyAlreadyExists
	}

	return nil
}

func (i *idempotencyRepository) Complete(ctx context.Context, idempotency domain.Idempotency) error {
	q := querier(ctx, i.pool)

	commandTag, err := q.Exec(ctx, `update idempotency_keys set transaction_id = $1, status = $2, response = $3,updated_at = $4 where idempotency_key = $5 and user_id = $6`,
		idempotency.TransactionId, idempotency.Status, idempotency.Response, idempotency.UpdatedAt, idempotency.IdempotencyKey, idempotency.UserId,
	)

	if err != nil {
		return fmt.Errorf("Error in completing of idempotency key: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrorNoRowsChanged
	}

	return nil
}
