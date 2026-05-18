package repository

import (
	"context"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/pkg/pagination"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{
		pool: pool,
	}
}

func (tr *transactionRepository) GetByUserId(ctx context.Context, userId int, cursor *pagination.TransactionCursor, limit int, currencies *[]string) ([]domain.Transaction, error) {
	query := `
		select tr.id, tr.type, tr.amount, tr.account_id, tr.status, tr.status_message, tr.created_at, tr.currency_id, 
		c.code as currency_code, c.symbol as currency_symbol, 



		d.source as deposit_source, 
		w.source as withdrawals_source,

		ua.firstname as sender_first_name,
		ua.lastname as sender_last_name,
		ua.phone_number as sender_phone_number,


		ub.firstname as receiver_first_name,
		ub.lastname as receiver_last_name,
		ub.phone_number as receiver_phone_number



		from transactions tr join currencies c on tr.currency_id = c.id

		join accounts aa on aa.id = tr.account_id

		left join deposits d on d.transaction_id = tr.id
		left join withdrawals w on w.transaction_id = tr.id
		left join transfers t on t.transaction_id = tr.id

		left join accounts a on a.id = t.sender_account_id
		left join accounts b on b.id = t.receiver_account_id

		left join users ua on ua.id = a.user_id
		left join users ub on ub.id = b.user_id


		where aa.user_id = $1
	`

	args := []any{userId}
	argsCount := 2

	if currencies != nil {
		query += fmt.Sprintf(` and c.code = ANY($%d)`, argsCount)
		argsCount += 1
		args = append(args, *currencies)
	}

	if cursor != nil {
		query += fmt.Sprintf(` and (tr.created_at, tr.id) < ($%d, $%d)`, argsCount, argsCount+1)
		args = append(args, cursor.CreatedAt, cursor.Id)
		argsCount += 2
	}

	query += fmt.Sprintf(` order by tr.created_at desc, tr.id desc limit $%d`, argsCount)

	args = append(args, limit)

	rows, err := tr.pool.Query(ctx, query, args...)

	if err != nil {
		return []domain.Transaction{}, fmt.Errorf("Error in rows: %w", err)
	}
	defer rows.Close()

	sl := make([]domain.Transaction, 0)

	for rows.Next() {
		var depositSource *string
		var withdrawalSource *string

		var senderFirstName *string
		var senderLastName *string
		var senderPhoneNumber *string

		var receiverFirstName *string
		var receiverLastName *string
		var receiverPhoneNumber *string

		var transaction domain.Transaction

		if err := rows.Scan(&transaction.Id, &transaction.Type, &transaction.Amount, &transaction.AccountId, &transaction.Status, &transaction.StatusMessage,
			&transaction.CreatedAt, &transaction.CurrencyId, &transaction.CurrencyCode, &transaction.CurrencySymbol,
			&depositSource, &withdrawalSource,

			&senderFirstName, &senderLastName, &senderPhoneNumber,
			&receiverFirstName, &receiverLastName, &receiverPhoneNumber,
		); err != nil {
			return []domain.Transaction{}, fmt.Errorf("Error in scanning: %w", err)
		}

		if depositSource != nil {
			transaction.DepositDetail = &domain.DepositDetail{
				Source: *depositSource,
			}
		}

		if withdrawalSource != nil {
			transaction.WithDrawalDetail = &domain.WithdrawalDetail{
				Source: *withdrawalSource,
			}
		}

		if senderFirstName != nil {
			transaction.TransferDetail = &domain.TransferDetail{
				Sender: domain.UserDetail{
					FirstName:   *senderFirstName,
					LastName:    *senderLastName,
					PhoneNumber: *senderPhoneNumber,
				},
				Receiver: domain.UserDetail{
					FirstName:   *receiverFirstName,
					LastName:    *receiverLastName,
					PhoneNumber: *receiverPhoneNumber,
				},
			}
		}

		sl = append(sl, transaction)

	}

	if err := rows.Err(); err != nil {
		return []domain.Transaction{}, fmt.Errorf("Error in row: %w", err)
	}

	return sl, nil

}

func (tr *transactionRepository) Create(ctx context.Context, t ...domain.Transaction) (map[int]int, error) {
	q := querier(ctx, tr.pool)

	mp := make(map[int]int, len(t))

	for _, val := range t {
		var transactionId int
		err := q.QueryRow(ctx, `insert into 
	transactions(type, amount, account_id, status_message, currency_id) 
	values ($1,$2,$3,$4,$5) returning id`, val.Type, val.Amount, val.AccountId, val.StatusMessage, val.CurrencyId).Scan(&transactionId)

		if err != nil {
			return map[int]int{}, fmt.Errorf("Error occured while creating new transfer: %w", err)
		}

		mp[val.AccountId] = transactionId

	}

	return mp, nil
}

func (tr *transactionRepository) MarkTransaction(ctx context.Context, status string, status_message string, id int) error {
	q := querier(ctx, tr.pool)

	commandTag, err := q.Exec(ctx, `update transactions set status = $1, status_message = $2 where id = $3`, status, status_message, id)

	if err != nil {
		return fmt.Errorf("Error in mark complete the transaction: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrorTransactionNotFound
	}

	return nil

}

// create table if not exists transactions(
//     id BIGINT generated always as identity PRIMARY KEY,
//     type transaction_type not null,
//     amount numeric(12,2) not null check (amount >0),
//     account_id BIGINT not null REFERENCES accounts(id),
//     status transaction_status not null DEFAULT 'pending',
//     status_message text not null,
//     created_at TIMESTAMPtz not null DEFAULT now()
// );

func (tr *transactionRepository) GetByAccountId(ctx context.Context, id int, limit int, transactionCursor *pagination.TransactionCursor) ([]domain.Transaction, int, error) {
	query := `select 
		tr.id, tr.type, tr.amount, tr.account_id, tr.currency_id,tr.status, 
		tr.status_message, tr.created_at, 

		d.source as deposit_source, 

		w.source as withdrawal_source, 

		c.code as currency_code,
		c.symbol as currency_symbol,

		au.id as user_id,

		au.firstname as sender_firstname,
		au.lastname as sender_lastname,
		au.phone_number as sender_phone_number,

		bu.firstname as receiver_firstname,
		bu.lastname as receiver_lastname,
		bu.phone_number as receiver_phone_number

		from transactions tr 
		left join deposits d on d.transaction_id = tr.id 
		left join withdrawals w on w.transaction_id = tr.id 
		left join transfers t on t.transaction_id = tr.id
		left join accounts a on a.id = t.sender_account_id
		left join accounts b on b.id = t.receiver_account_id
		left join users au on au.id = a.user_id
		left join users bu on bu.id = b.user_id
		left join currencies c on c.id = tr.currency_id

		where tr.account_id = $1
	`

	args := []any{id}

	argsCount := 2

	if transactionCursor != nil {
		query += fmt.Sprintf(" and (tr.created_at,tr.id) < ($%d, $%d)", argsCount, argsCount+1)
		argsCount += 2
		args = append(args, transactionCursor.Id, transactionCursor.CreatedAt)
	}

	query += fmt.Sprintf(` order by created_at desc, tr.id desc limit $%d`, argsCount)

	args = append(args, limit)

	rows, err := tr.pool.Query(ctx,
		query, args...,
	)

	if err != nil {
		return []domain.Transaction{}, 0, fmt.Errorf("Error in get transactions by id here: %w", err)
	}

	sl := make([]domain.Transaction, 0)
	var userId int

	for rows.Next() {
		var t domain.Transaction
		var (
			withdrawalDetailSource *string
			depositDetailSource    *string

			senderFirstName   *string
			senderLastName    *string
			senderPhoneNumber *string

			receiverFirstName   *string
			receiverLastName    *string
			receiverPhoneNumber *string
		)

		if err := rows.Scan(&t.Id, &t.Type, &t.Amount, &t.AccountId, &t.CurrencyId, &t.Status, &t.StatusMessage, &t.CreatedAt,
			&depositDetailSource, &withdrawalDetailSource,
			&t.CurrencyCode, &t.CurrencySymbol,
			&userId,
			&senderFirstName, &senderLastName, &senderPhoneNumber,
			&receiverFirstName, &receiverLastName, &receiverPhoneNumber,
		); err != nil {
			return []domain.Transaction{}, 0, fmt.Errorf("Error in a loop get transactions by id: %w", err)
		}
		// 'transfer', 'deposit','withdraw'

		switch t.Type {
		case "transfer":
			if senderFirstName != nil {
				t.TransferDetail = &domain.TransferDetail{
					Sender: domain.UserDetail{
						FirstName:   *senderFirstName,
						LastName:    *senderLastName,
						PhoneNumber: *senderPhoneNumber,
					},
					Receiver: domain.UserDetail{
						FirstName:   *receiverFirstName,
						LastName:    *receiverLastName,
						PhoneNumber: *receiverPhoneNumber,
					},
				}
			}
		case "deposit":
			if depositDetailSource != nil {
				t.DepositDetail = &domain.DepositDetail{
					Source: *depositDetailSource,
				}
			}

		case "withdraw":
			if withdrawalDetailSource != nil {
				t.WithDrawalDetail = &domain.WithdrawalDetail{
					Source: *withdrawalDetailSource,
				}
			}
		}
		sl = append(sl, t)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return []domain.Transaction{}, 0, fmt.Errorf("Error in a row get transactions by id: %w", err)
	}

	return sl, userId, nil
}

func (tr *transactionRepository) GetAll(ctx context.Context) ([]domain.Transaction, error) {

	rows, err := tr.pool.Query(ctx,
		`select 
		tr.id, tr.type, tr.amount, tr.account_id, tr.currency_id,tr.status, 
		tr.status_message, tr.created_at, 

		d.source as deposit_source, 

		w.source as withdrawal_source, 

		c.code as currency_code,
		c.symbol as currency_symbol,

		au.firstname as sender_firstname,
		au.lastname as sender_lastname,
		au.phone_number as sender_phone_number,

		bu.firstname as receiver_firstname,
		bu.lastname as receiver_lastname,
		bu.phone_number as receiver_phone_number

		from transactions tr 
		left join deposits d on d.transaction_id = tr.id 
		left join withdrawals w on w.transaction_id = tr.id 
		left join transfers t on t.transaction_id = tr.id
		left join accounts a on a.id = t.sender_account_id
		left join accounts b on b.id = t.receiver_account_id
		left join users au on au.id = a.user_id
		left join users bu on bu.id = b.user_id
		left join currencies c on c.id = tr.currency_id
	`)

	if err != nil {
		return []domain.Transaction{}, fmt.Errorf("Error in get transactions by id sdfdsf: %w", err)
	}

	sl := make([]domain.Transaction, 0)

	for rows.Next() {
		var t domain.Transaction
		var (
			withdrawalDetailSource *string
			depositDetailSource    *string

			senderFirstName   *string
			senderLastName    *string
			senderPhoneNumber *string

			receiverFirstName   *string
			receiverLastName    *string
			receiverPhoneNumber *string
		)

		if err := rows.Scan(&t.Id, &t.Type, &t.Amount, &t.AccountId, &t.CurrencyId, &t.Status, &t.StatusMessage, &t.CreatedAt,
			&withdrawalDetailSource, &depositDetailSource,
			&t.CurrencyCode, &t.CurrencySymbol,
			&senderFirstName, &senderLastName, &senderPhoneNumber,
			&receiverFirstName, &receiverLastName, &receiverPhoneNumber,
		); err != nil {
			return []domain.Transaction{}, fmt.Errorf("Error in a loop get transactions by id: %w", err)
		}
		// 'transfer', 'deposit','withdraw'

		switch t.Type {
		case "transfer":
			if senderFirstName != nil {
				t.TransferDetail = &domain.TransferDetail{
					Sender: domain.UserDetail{
						FirstName:   *senderFirstName,
						LastName:    *senderLastName,
						PhoneNumber: *senderPhoneNumber,
					},
					Receiver: domain.UserDetail{
						FirstName:   *receiverFirstName,
						LastName:    *receiverLastName,
						PhoneNumber: *receiverPhoneNumber,
					},
				}
			}
		case "deposit":
			if depositDetailSource != nil {
				t.DepositDetail = &domain.DepositDetail{
					Source: *depositDetailSource,
				}
			}

		case "withdraw":
			if withdrawalDetailSource != nil {
				t.WithDrawalDetail = &domain.WithdrawalDetail{
					Source: *withdrawalDetailSource,
				}
			}
		}
		sl = append(sl, t)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return []domain.Transaction{}, fmt.Errorf("Error in a row get transactions by id: %w", err)
	}

	return sl, nil
}

// CREATE TABLE IF NOT EXISTS transfers(
//     id BIGINT generated always as identity PRIMARY key,
//     transaction_id BIGINT REFERENCES transactions(id),
//     sender_account_id BIGINT not null REFERENCES accounts(id),
//     receiver_account_id BIGINT not null REFERENCES accounts(id),
//     currency_id BIGINT not null REFERENCES currencies(id),
//     amount numeric(12,2) not null check (amount > 0)
// );
// create table if not exists deposits(
//     id BIGINT generated always as identity PRIMARY KEY,
//     transaction_id BIGINT REFERENCES transactions(id) not null,
//     account_id BIGINT REFERENCES accounts(id) not null,
//     amount numeric(12,2) not null check (amount>0),
//     source deposit_source not null
// );
// create table if not exists withdrawals(
//     id BIGINT generated always as identity PRIMARY KEY,
//     transaction_id BIGINT REFERENCES transactions(id) not null,
//     account_id BIGINT REFERENCES accounts(id) not NULL,
//     amount numeric(12,2) not null check (amount > 0),
//     source withdraw_source not null
// );
// CREATE TABLE IF NOT EXISTS USERS (
//     id BIGINT generated always as identity primary key,
//     firstname text not null,
//     lastname text not null,
//     birthday TIMESTAMPtz not null,
//     phone_number text not null,
//     password text not null,
//     created_at TIMESTAMPtz not null default now()
// );

// CREATE TABLE IF NOT EXISTS CURRENCIES(
//     id bigint generated always as identity primary key,
//     name text not null,
//     code char(3) not null unique,
//     symbol VARCHAR(5) not null,
//     created_at TIMESTAMPtz not null DEFAULT now()
// );

// create table accounts(
//     id BIGINT generated always as identity PRIMARY KEY,
//     user_id BIGINT not null REFERENCES users(id),
//     currency_id BIGINT not null REFERENCES currencies(id),
//     balance numeric(12,2) not null DEFAULT 0,
//     is_active BOOLEAN not null DEFAULT true,
//     created_at TIMESTAMPTZ not null DEFAULT now()
// );

// create table if not exists transactions(
//     id BIGINT generated always as identity PRIMARY KEY,
//     type transaction_type not null,
//     amount numeric(12,2) not null check (amount >0),
//     account_id BIGINT not null REFERENCES accounts(id),
//     status transaction_status not null DEFAULT 'pending',
//     status_message text not null,
//     created_at TIMESTAMPtz not null DEFAULT now()
// );
