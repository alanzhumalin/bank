package repository

import (
	"context"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type accountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) AccountRepository {
	return &accountRepository{pool: pool}
}

func (a *accountRepository) Create(ctx context.Context, acc domain.Account) error {
	_, err := a.pool.Exec(ctx, `insert into accounts(user_id, currency_id) values($1,$2)`, acc.UserId, acc.CurrencyId)

	if err != nil {
		return fmt.Errorf("Error in creating account, account_repository: %w", err)
	}

	return nil
}

func (a *accountRepository) DeleteById(ctx context.Context, id int) error {
	tag, err := a.pool.Exec(ctx, `delete from accounts where id = $1`, id)

	if err != nil {
		return fmt.Errorf("Error in deleting by id the account: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.AccountNotFound
	}

	return nil
}

func (a *accountRepository) GetAll(ctx context.Context) ([]domain.Account, error) {
	rows, err := a.pool.Query(ctx, `select id, user_id, currency_id, balance,is_active, created_at from accouns`)

	if err != nil {
		return []domain.Account{}, fmt.Errorf("Error in getting all accounts: %w", err)
	}

	sl := make([]domain.Account, 0)
	for rows.Next() {
		var acc domain.Account

		if err := rows.Scan(&acc.Id, &acc.UserId, &acc.CurrencyId, &acc.Balance, &acc.IsActive, &acc.CreatedAt); err != nil {
			return []domain.Account{}, fmt.Errorf("Error in getting all accounts, in a loop: %w", err)
		}

		sl = append(sl, acc)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return []domain.Account{}, fmt.Errorf("Error in getting all accounts, in structure rows: %w", err)
	}

	return sl, nil
}
