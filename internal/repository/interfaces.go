package repository

import (
	"context"

	"github.com/alanzhumalin/bank/internal/domain"
	user "github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
)

type UserRepository interface {
	Create(ctx context.Context, u user.User) error
	UserExists(ctx context.Context, phoneNumber string) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user domain.User) error
	GetByPhone(ctx context.Context, phone string) (domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
}

type CurrencyRepository interface {
	Create(ctx context.Context, c domain.Сurrency) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (domain.Сurrency, error)
	UpdateById(ctx context.Context, id int, name string, code string, symbol string) error
	Exists(ctx context.Context, code string) (bool, error)
	GetAll(ctx context.Context) ([]domain.Сurrency, error)
}

type TransferRepository interface {
}

type AccountRepository interface {
	Create(ctx context.Context, a domain.Account) error
	DeleteById(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]domain.Account, error)
	SelectTwoAccountsForUpdate(ctx context.Context, senderAccountId int, receiverAccountId int) (domain.Account, domain.Account, error)
	IncreaseBalance(ctx context.Context, balance decimal.Decimal, accountId int) error
	DecreaseBalance(ctx context.Context, balance decimal.Decimal, accountId int) error
}

type TxManagerRepository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TransactionRepository interface {
	Create(ctx context.Context, t domain.Transaction) (int, error)
	MarkTransaction(ctx context.Context, status string, status_message string, id int) error
}
