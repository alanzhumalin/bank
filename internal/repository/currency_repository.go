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
	err := c.pool.QueryRow(ctx, `select id, name, code, symbol, created_at from currencies`).Scan(&currency.Id, &currency.Name, &currency.Code, &currency.Symbol, &currency.CreatedAt)

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

func (c *currentRepository) Exists(ctx context.Context, code string) (bool, error) {
	var id int

	err := c.pool.QueryRow(ctx, `select id from currencies where code = $1`, code).Scan(&id)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	return false, fmt.Errorf("Error exists in currency_repository: %w", err)

}

func (c *currentRepository) GetAll(ctx context.Context) ([]domain.Сurrency, error) {
	var currencies []domain.Сurrency

	rows, err := c.pool.Query(ctx, `select id,name,code,symbol,created_at from currencies`)

	if err != nil {
		return []domain.Сurrency{}, fmt.Errorf("Error getall currency in current_repository: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var currency domain.Сurrency
		err := rows.Scan(&currency.Id, &currency.Name, &currency.Code, &currency.Symbol, &currency.CreatedAt)

		if err != nil {
			return []domain.Сurrency{}, fmt.Errorf("Error get row currency in current_repository: %w", err)
		}

		currencies = append(currencies, currency)
	}

	if err := rows.Err(); err != nil {
		return []domain.Сurrency{}, fmt.Errorf("Error rows in current_repository: %w", err)
	}
	return currencies, nil
}
