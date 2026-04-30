package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transferRepository struct {
	pool *pgxpool.Pool
}

type account struct {
	id          int
	currency_id int
	balance     float64
	is_active   bool
}

func NewTransferRepository(pool *pgxpool.Pool) TransferRepository {
	return &transferRepository{pool: pool}
}

func (tr *transferRepository) GetAll(ctx context.Context) ([]domain.Transfer, error) {
	transfers := make([]domain.Transfer, 0)

	rows, err := tr.pool.Query(ctx, `select tr.id, tr.sender_account_id, 
	tr.receiver_account_id, tr.currency_id, 
	tr.amount, tr.status, tr.status_message, tr.created_at

	c.name as currency_name, c.code as currency_code, c.symbol as currency_symbol, 
	
	sender.firstname as sender_firstname, sender.lastname as sender_lastname, 
	receiver.firstname as receiver_firstname, receiver.lastname as receiver_lastname, 
	
	from transfers tr join currencies c on tr.currency_id = c.id
	join accounts senderacc on senderacc.id = tr.sender_account_id
	join accounts receiveracc on receiveracc.id = tr.receiver_account_id
	join users sender on senderacc.user_id = sender.id
	join users receiver on receiveracc.user_id = receiver.id

	`)

	if err != nil {
		return []domain.Transfer{}, fmt.Errorf("Get all transfers, error in getting transfer: %w", err)
	}

	for rows.Next() {
		var transfer domain.Transfer

		err := rows.Scan(&transfer.Id, &transfer.SenderAccountId, &transfer.ReceiverAccountId, &transfer.CurrencyId, &transfer.Amount,
			&transfer.Status, &transfer.StatusMessage, &transfer.CreatedAt,
			&transfer.Currency.Name, &transfer.Currency.Code, &transfer.Currency.Symbol,
			&transfer.Sender.FirstName, &transfer.Sender.LastName,
			&transfer.Receiver.FirstName, &transfer.Receiver.LastName,
		)

		if err != nil {
			return []domain.Transfer{}, fmt.Errorf("Get all transfers, error in getting a transfer in a loop: %w", err)
		}

		transfers = append(transfers, transfer)

	}
	rows.Close()

	if err := rows.Err(); err != nil {
		return []domain.Transfer{}, fmt.Errorf("Get all transfers, error in row after loop: %w", err)
	}

	return transfers, nil

}

func (tr *transferRepository) GetById(ctx context.Context, id int) (domain.Transfer, error) {
	var transfer domain.Transfer

	err := tr.pool.QueryRow(ctx, `select tr.id, tr.sender_account_id, 
	tr.receiver_account_id, tr.currency_id, 
	tr.amount, tr.status, tr.status_message,
	tr.created_at

	c.name as currency_name, c.code as currency_code, c.symbol as currency_symbol, 
	
	sender.firstname as sender_firstname, sender.lastname as sender_lastname,

	receiver.firstname as receiver_firstname, receiver.lastname as receiver_lastname, receiver.phone_number as receiver_phone_number
	
	from transfers tr join currencies c on tr.currency_id = c.id
	join accounts senderacc on senderacc.id = tr.sender_account_id
	join accounts receiveracc on senderacc.id = tr.receiver_account_id
	join users sender on senderacc.user_id = sender.id
	join users receiver on receiveracc.user_id = receiver.id

	where tr.id = $1
	
	`, id).Scan(&transfer.Id, &transfer.SenderAccountId, &transfer.ReceiverAccountId, &transfer.CurrencyId, &transfer.Amount,
		&transfer.Status, &transfer.StatusMessage, &transfer.CreatedAt,
		&transfer.Currency.Name, &transfer.Currency.Code, &transfer.Currency.Symbol,
		&transfer.Sender.FirstName, &transfer.Sender.LastName,
		&transfer.Receiver.FirstName, &transfer.Receiver.LastName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Transfer{}, domain.ErrorTransferNotFound
		}
		return domain.Transfer{}, err
	}

	return transfer, nil

}

func (tr *transferRepository) Create(ctx context.Context, t domain.Transfer) error {
	tx, err := tr.pool.Begin(ctx)

	if err != nil {
		return fmt.Errorf("Create transfer, error in begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	var transactionId int

	err = tx.QueryRow(ctx, `insert into transfers(sender_account_id, receiver_account_id, currency_id, amount, status,status_message) values($1, $2,$3,$4,$5, $6) returning id`, t.SenderAccountId, t.ReceiverAccountId, t.CurrencyId, t.Amount, "pending", "created transfer").Scan(&transactionId)

	if err != nil {
		return fmt.Errorf("Create transfer, error in inserting transfer: %w", err)
	}
	rows, err := tx.Query(ctx, `select id, currency_id, balance, is_active from accounts where id = ANY($1) order by id asc for update`, []int{t.ReceiverAccountId, t.SenderAccountId})

	if err != nil {
		return fmt.Errorf("Create transfer, error in selecting account: %w", err)
	}

	defer rows.Close()

	mp := make(map[int]account)

	for rows.Next() {
		var account account
		if err := rows.Scan(&account.id, &account.currency_id, &account.balance, &account.is_active); err != nil {
			return fmt.Errorf("Create transfer, error in selecting account using next: %w", err)
		}
		mp[account.id] = account
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("Create transfer, error in after loop: %w", err)
	}

	receiverAccount, ok := mp[t.ReceiverAccountId]
	if !ok {
		return domain.AccountNotFound
	}
	senderAccount, ok := mp[t.SenderAccountId]
	if !ok {
		return domain.AccountNotFound
	}

	if receiverAccount.currency_id != t.CurrencyId {
		return domain.AccountNotSupportCurrency
	}

	if !receiverAccount.is_active {
		return domain.AccountIsNotActive
	}

	if senderAccount.balance < t.Amount {
		return domain.ErrorNotEnoughBalance
	}

	tag, err := tx.Exec(ctx, `update accounts set balance = balance - $1 where id = $2`, t.Amount, t.SenderAccountId)

	if err != nil {
		return fmt.Errorf("Create transfer, error in updating the balance of the sender account: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.AccountNotFound
	}

	tag, err = tx.Exec(ctx, `update accounts set balance = balance + $1 where id = $2`, t.Amount, t.ReceiverAccountId)

	if err != nil {
		return fmt.Errorf("Create transfer, error in updating the balance of the receiver account: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.AccountNotFound
	}

	tag, err = tx.Exec(ctx, `update transfers set status = $1, status_message = $2 where id = $3`, "completed", "the transfer is completed", transactionId)

	if err != nil {
		return fmt.Errorf("Create transfer, error in updating the status of the transfer: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrorTransferNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("Create transfer, error in commiting the transaction: %w", err)
	}

	return nil

}
