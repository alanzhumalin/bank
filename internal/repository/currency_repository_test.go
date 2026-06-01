package repository

import (
	"context"
	"testing"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func RunDbContainer(ctx context.Context, t *testing.T) *pgxpool.Pool {
	pgContainer, err := postgres.Run(
		ctx, "postgres:16-alpine",
		postgres.BasicWaitStrategies(),
		postgres.WithDatabase("bank-test"),
		postgres.WithUsername("admin"),
		postgres.WithPassword("admin"),
	)

	if err != nil {
		t.Fatalf("run pg container err: %v", err)
	}

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		t.Fatalf("connection string get error: %v", err)
	}

	pool, err := pgxpool.New(ctx, connString)

	if err != nil {
		t.Fatalf("create new pool error: %v", err)
	}

	return pool
}

func InitDb(ctx context.Context, t *testing.T, pool *pgxpool.Pool) {
	_, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS CURRENCIES(
		id bigint generated always as identity primary key,
		name text not null,
		code char(3) not null unique, 
		symbol VARCHAR(5) not null,
		created_at TIMESTAMPtz not null DEFAULT now()
	);`)

	if err != nil {
		t.Fatalf("creating table currencies error: %v", err)
	}
}

func TestCurrencyRepositoryCreateCurrency(t *testing.T) {
	ctx := context.Background()
	pool := RunDbContainer(ctx, t)

	InitDb(ctx, t, pool)

	repository := NewCurrencyRepository(pool)

	newC := domain.Сurrency{
		Name:   "Dollar",
		Code:   "USD",
		Symbol: "$",
	}

	err := repository.Create(ctx, newC)

	if err != nil {
		t.Fatalf("create new currency error: %v", err)
	}

	c, err := repository.GetById(ctx, 1)

	if err != nil {
		t.Fatalf("get currency by id error: %v", err)
	}

	if c.Name != newC.Name || c.Code != newC.Code || c.Symbol != newC.Symbol {
		t.Fatal("expected to be equal")
	}
}
