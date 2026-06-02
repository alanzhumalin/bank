package dto

import "github.com/shopspring/decimal"

// create table if not exists withdrawals(
//     id BIGINT generated always as identity PRIMARY KEY,
//     transaction_id BIGINT REFERENCES transactions(id) not null,
//     account_id BIGINT REFERENCES accounts(id) not NULL,
//     amount numeric(12,2) not null check (amount > 0),
//     source withdraw_source not null
// );

type CreateWindrawalRequest struct {
	AccountId      int             `json:"account_id"`
	Amount         decimal.Decimal `json:"amount"`
	Source         string          `json:"source"`
	IdempotencyKey string          `json:"idempotency_key"`
}

func (c *CreateWindrawalRequest) Validate() error {
	if c.AccountId <= 0 {
		return ErrorAccountIdRequired
	}
	if c.Amount.LessThanOrEqual(decimal.Zero) {
		return ErrorAmountRequired
	}

	if c.Source == "" {
		return ErrorSourceRequired
	}

	if c.IdempotencyKey == "" {
		return ErrorIdempotencyKeyIsRequired
	}

	return nil

}
