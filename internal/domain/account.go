package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// create table accounts(
//     id BIGINT generated always as identity PRIMARY KEY,
//     user_id BIGINT not null REFERENCES users(id),
//     currency_id BIGINT not null REFERENCES currencies(id),
//     balance numeric(12,2) not null DEFAULT 0,
//     is_active BOOLEAN not null DEFAULT true,
//     created_at TIMESTAMPTZ not null DEFAULT now()
// );

type Account struct {
	Id         int
	UserId     int
	CurrencyId int
	Balance    decimal.Decimal
	IsActive   bool
	CreatedAt  time.Time
}
