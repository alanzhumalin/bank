package repository

import (
	"context"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
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

func (tr *transactionRepository) Create(ctx context.Context, t ...domain.Transaction) (map[int]int, error) {
	q := querier(ctx, tr.pool)

	mp := make(map[int]int, len(t))

	for _, val := range t {
		var transactionId int
		err := q.QueryRow(ctx, `insert into 
	transactions(type, amount, account_id, status, status_message, currency_id) 
	values ($1,$2,$3,$4,$5, $6) returning id`, val.Type, val.Amount, val.AccountId, val.Status, val.StatusMessage, val.CurrencyId).Scan(&transactionId)

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

func (tr *transactionRepository) GetByAccountId(ctx context.Context, id int) ([]domain.Transaction, error) {
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

		where tr.account_id = $1

	`, id)

	if err != nil {
		return []domain.Transaction{}, fmt.Errorf("Error in get transactions by id: %w", err)
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
		return []domain.Transaction{}, fmt.Errorf("Error in get transactions by id: %w", err)
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
