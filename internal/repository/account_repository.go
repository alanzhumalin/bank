package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
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

func (a *accountRepository) SelectTwoAccountsForUpdate(ctx context.Context, senderAccountId int, receiverAccountId int) (domain.Account, domain.Account, error) {
	q := querier(ctx, a.pool)

	rows, err := q.Query(ctx, `select id, currency_id, balance, is_active from accounts where id = ANY($1) order by id for update`, []int{senderAccountId, receiverAccountId})

	if err != nil {
		return domain.Account{}, domain.Account{}, fmt.Errorf("Error in selecting two account for update: %w", err)
	}

	mp := make(map[int]domain.Account, 2)

	for rows.Next() {
		var req domain.Account

		err := rows.Scan(&req.Id, &req.CurrencyId, &req.Balance, &req.IsActive)
		if err != nil {
			return domain.Account{}, domain.Account{}, fmt.Errorf("Error occured in a rows loop: %w", err)
		}
		mp[req.Id] = req
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return domain.Account{}, domain.Account{}, fmt.Errorf("Error occured in a rows: %w", err)
	}

	senderAccount, ok := mp[senderAccountId]
	if !ok {
		return domain.Account{}, domain.Account{}, domain.AccountNotFound
	}

	receiverAccount, ok := mp[receiverAccountId]

	if !ok {
		return domain.Account{}, domain.Account{}, domain.AccountNotFound
	}

	return senderAccount, receiverAccount, nil

}

func (a *accountRepository) GetByIdForUpdate(ctx context.Context, id int) (domain.Account, error) {
	q := querier(ctx, a.pool)

	var account domain.Account

	err := q.QueryRow(ctx, `select id, user_id, currency_id, balance, is_active, created_at from accounts where id = $1 for update`, id).Scan(&account.Id, &account.UserId, &account.CurrencyId, &account.Balance, &account.IsActive)

	if err == nil {
		return account, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Account{}, domain.AccountNotFound
	}

	return domain.Account{}, fmt.Errorf("Error in get by id account: %w", err)

}

func (a *accountRepository) IncreaseBalance(ctx context.Context, balance decimal.Decimal, accountId int) error {
	q := querier(ctx, a.pool)

	commandTag, err := q.Exec(ctx, `update accounts set balance = balance + $1 where id = $2`, balance, accountId)

	if err != nil {
		return fmt.Errorf("Error in increasing the balance: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return domain.AccountNotFound
	}
	return nil
}

func (a *accountRepository) DecreaseBalance(ctx context.Context, balance decimal.Decimal, accountId int) error {
	q := querier(ctx, a.pool)

	commandTag, err := q.Exec(ctx, `update accounts set balance = balance - $1 where id = $2`, balance, accountId)

	if err != nil {
		return fmt.Errorf("Error in descreasing the balance: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return domain.AccountNotFound
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
