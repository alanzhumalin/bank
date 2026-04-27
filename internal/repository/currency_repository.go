package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type currentRepository struct {
	pool *pgxpool.Pool
}

func NewCurrencyRepository(pool *pgxpool.Pool) CurrencyRepository {
	return &currentRepository{pool: pool}
}

func (c *currentRepository) Delete(ctx context.Context, id int) error {
	res, err := c.pool.Exec(ctx, `delete from currencies where id = $1`, id)

	if err != nil {
		return fmt.Errorf("Error delete currency in currency_repository: %w", err)
	}

	if res.RowsAffected() == 0 {
		return domain.ErrorCurrencyNotFound
	}

	return nil
}
func (c *currentRepository) Create(ctx context.Context, currency domain.Сurrency) error {
	_, err := c.pool.Exec(ctx, `insert into currencies(name,code,symbol) values ($1, $2, $3)`, currency.Name, currency.Code, currency.Symbol)

	if err != nil {
		return fmt.Errorf("Error create currency in currency_repository: %w", err)
	}

	return nil
}
func (c *currentRepository) GetById(ctx context.Context, id int) (domain.Сurrency, error) {
	var currency domain.Сurrency
	err := c.pool.QueryRow(ctx, `select * from currencies`).Scan(&currency.Id, &currency.Name, &currency.Code, &currency.Symbol, &currency.Created_at)

	if err == nil {
		return currency, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Сurrency{}, domain.ErrorCurrencyNotFound
	}

	return domain.Сurrency{}, fmt.Errorf("Error get by id in currency_repository: %w", err)
}
func (c *currentRepository) UpdateById(ctx context.Context, id int, name string, code string, symbol string) error {
	res, err := c.pool.Exec(ctx, `update currencies set name = $1, code = $2, symbol = $3 where id = $4`, name, code, symbol, id)

	if err != nil {
		return fmt.Errorf("Error update currency by id in currency_repository: %w", err)
	}

	if res.RowsAffected() == 0 {
		return domain.ErrorCurrencyNotFound
	}

	return nil
}
