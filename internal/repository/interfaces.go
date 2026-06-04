package repository

import (
	"context"
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	user "github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/pkg/pagination"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
)

type UserRepository interface {
	Create(ctx context.Context, u user.User) (int, string, error)
	UserExists(ctx context.Context, phoneNumber string) (bool, error)
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
	Create(ctx context.Context, t ...domain.Transfer) error
}

type AccountRepository interface {
	Create(ctx context.Context, a domain.Account) error
	DeleteById(ctx context.Context, id int, time time.Time) error
	GetAll(ctx context.Context) ([]domain.Account, error)
	SelectTwoAccountsForUpdate(ctx context.Context, senderAccountId int, receiverAccountId int) (domain.Account, domain.Account, error)
	IncreaseBalance(ctx context.Context, balance decimal.Decimal, accountId int) error
	DecreaseBalance(ctx context.Context, balance decimal.Decimal, accountId int) error
	GetByIdForUpdate(ctx context.Context, id int) (domain.Account, error)
	GetUserAccounts(ctx context.Context, userId int) ([]domain.Account, error)
	Exists(ctx context.Context, userId int, currencyId int) (bool, error)
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
	Create(ctx context.Context, t ...domain.Transaction) (map[int]int, error)
	MarkTransaction(ctx context.Context, status string, status_message string, id int) error
	GetByAccountId(ctx context.Context, id int, limit int, transactionCursor *pagination.TransactionCursor) ([]domain.Transaction, int, error)
	GetAll(ctx context.Context) ([]domain.Transaction, error)
	GetByUserId(ctx context.Context, userId int, cursor *pagination.TransactionCursor, limit int, currencies *[]string) ([]domain.Transaction, error)
}

type DepositRepository interface {
	Create(ctx context.Context, d domain.Deposit) error
}

type WithdrawalRepository interface {
	Create(ctx context.Context, w domain.Withdrawal) error
}

type AuthRepository interface {
	GetDetails(context context.Context, phoneNumber string) (LoginDetails, error)
	Сreate(ctx context.Context, session domain.Session) error
	Revoke(ctx context.Context, sessionId string) error
	RevokeAllUserDevices(ctx context.Context, id int) error
	Update(ctx context.Context, newHashedToken string, expires_at time.Time, sessionId string) error
	GetSessionById(ctx context.Context, sessionId string) (domain.Session, error)
}

type IdempotencyRepository interface {
	GetByKey(ctx context.Context, key string, userId int) (domain.Idempotency, error)
	Start(ctx context.Context, idempotency domain.Idempotency) error
	Complete(ctx context.Context, idempotency domain.Idempotency) error
}
