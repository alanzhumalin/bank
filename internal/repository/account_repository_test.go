package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func runDb(t *testing.T, ctx context.Context) *pgxpool.Pool {

	t.Helper()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("bank_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		t.Fatalf("Error occured %v", err)
	}

	t.Cleanup(func() {
		_ = pgContainer.Terminate(ctx)
	})

	conn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		t.Fatalf("Error occured %v", err)
	}

	pool, err := pgxpool.New(ctx, conn)

	if err != nil {
		t.Fatalf("Error occured %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool

}

func initializeDb(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
	t.Helper()

	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS CURRENCIES(
		id bigint generated always as identity primary key,
		name text not null,
		code char(3) not null unique, 
		symbol VARCHAR(5) not null,
		created_at TIMESTAMPtz not null DEFAULT now()
	);
	`)

	if err != nil {
		t.Fatalf("Error occured %v", err)
	}

	_, err = pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS USERS (
			id BIGINT generated always as identity primary key,
			firstname text not null,
			lastname text not null,
			birthday TIMESTAMPtz not null,
			phone_number text not null,
			password text not null,
			created_at TIMESTAMPtz not null default now()
		);
	`)

	if err != nil {
		t.Fatalf("Error occured %v", err)
	}

	_, err = pool.Exec(ctx, `create table accounts(
		id BIGINT generated always as identity PRIMARY KEY,
		user_id BIGINT not null REFERENCES users(id),
		currency_id BIGINT not null REFERENCES currencies(id),
		balance numeric(12,2) not null DEFAULT 0,
		is_active BOOLEAN not null DEFAULT true,
		created_at TIMESTAMPTZ not null DEFAULT now()
	);
	`)

	if err != nil {
		t.Fatalf("Error occured %v", err)
	}

	commandTag, err := pool.Exec(ctx, `insert into currencies(name,code,symbol) values('dollar', 'USD', '$')`)

	if err != nil {
		t.Fatalf("Error occured %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		t.Fatalf("No rows inserted")
	}

	commandTag, err = pool.Exec(ctx, `insert into users(firstname,lastname, birthday, phone_number, password) 
		values('Alan', 'Walker', '2006-01-02', '87781234112', 'Swomehtgin10006')
	`)

	if err != nil {
		t.Fatalf("Error occured %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		t.Fatalf("No rows inserted")
	}

	commandTag, err = pool.Exec(ctx, `insert into users(firstname,lastname, birthday, phone_number, password) 
		values('Zack', 'Walker', '2006-01-02', '87051314111', 'Something10002')
	`)

	if err != nil {
		t.Fatalf("Error occured %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		t.Fatalf("No rows inserted")
	}

	commandTag, err = pool.Exec(ctx, `insert into accounts(user_id, currency_id, balance) values(1, 1, 100000)`)

	if err != nil {
		t.Fatalf("Error occured %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		t.Fatalf("No rows inserted")
	}

	commandTag, err = pool.Exec(ctx, `insert into accounts(user_id, currency_id, balance) values(2, 1, 100000)`)

	if err != nil {
		t.Fatalf("Error occured %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		t.Fatalf("No rows inserted")
	}
}

func TestAccountRepositorySelectTwoForUpdate(t *testing.T) {
	ctx := context.Background()
	pool := runDb(t, ctx)

	initializeDb(t, pool, ctx)

	accountRepository := NewAccountRepository(pool)

	acc1, acc2, err := accountRepository.SelectTwoAccountsForUpdate(ctx, 1, 2)

	if err != nil {
		t.Fatalf("Error occured %v", err)
	}

	if acc1.Id != 1 {
		t.Fatalf("Expected id %v, got %v", 1, acc1.Id)
	}

	if acc2.Id != 2 {
		t.Fatalf("Expected id %v, got %v", 2, acc2.Id)
	}

	if !acc1.IsActive {
		t.Fatalf("Expected acc1 to be active")
	}

	if !acc2.IsActive {
		t.Fatalf("Expected acc2 to be active")
	}

	if acc1.CurrencyId != 1 {
		t.Fatalf("Expected currency id of acc1 to be 1")
	}

	if acc2.CurrencyId != 1 {
		t.Fatalf("Expected currency id of acc2 to be 1")
	}

}

func TestAccountRepositorySelectTwoForUpdateSenderNotFound(t *testing.T) {
	ctx := context.Background()
	pool := runDb(t, ctx)

	initializeDb(t, pool, ctx)

	accountRepository := NewAccountRepository(pool)

	_, _, err := accountRepository.SelectTwoAccountsForUpdate(ctx, 3, 2)

	if !errors.Is(err, domain.AccountNotFound) {
		t.Fatalf("Expected account not found")
	}

}

func TestAccountRepositorySelectTwoForUpdateReceiverNotFound(t *testing.T) {
	ctx := context.Background()

	pool := runDb(t, ctx)

	initializeDb(t, pool, ctx)

	accountRepository := NewAccountRepository(pool)

	_, _, err := accountRepository.SelectTwoAccountsForUpdate(ctx, 1, 3)

	if !errors.Is(err, domain.AccountNotFound) {
		t.Fatalf("Expected account not found")
	}

}

func TestAccountRepositorySelectTwoForUpdateWithTx(t *testing.T) {
	ctx := context.Background()

	pool := runDb(t, ctx)

	initializeDb(t, pool, ctx)

	accountRepo := NewAccountRepository(pool)

	tx, err := pool.Begin(ctx)

	if err != nil {
		t.Fatalf("Error occured in transaction: %v", err)
	}

	defer tx.Rollback(ctx)

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	_, _, err = accountRepo.SelectTwoAccountsForUpdate(ctxWithTx, 1, 2)

	if err != nil {
		t.Fatalf("Error occured: %v", err)
	}

	done := make(chan error, 1)

	go func() {

		ctxWithTimer, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		defer cancel()

		_, err := pool.Exec(ctxWithTimer, `update accounts set balance = balance + 1 where id = ANY($1)`, []int{1, 2})

		done <- err

	}()

	err = <-done

	if err == nil {
		t.Fatal("Expected to be non nil error")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Expected to be error: %v, got %v", context.DeadlineExceeded, err)
	}

}
