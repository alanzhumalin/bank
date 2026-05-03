package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func querier(ctx context.Context, pool *pgxpool.Pool) Querier {
	tx, ok := GetTx(ctx)
	if !ok {
		return pool
	}
	return tx
}
