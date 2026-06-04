package domain

import "time"

// create table if not exists idempotency_keys(
//     id BIGINT generated always as identity PRIMARY KEY,
//     transaction_id BIGINT REFERENCES transactions(id),
//     user_id BIGINT REFERENCES users(id) not null,
//     idempotency_key text not null,
//     operation text not null,
//     status text not null DEFAULT 'pending',
//     response JSONB,
//     created_at TIMESTAMPTZ not null DEFAULT now(),
//     updated_at TIMESTAMPTZ not null,
//     UNIQUE(user_id, idempotency_key)
// );

type Idempotency struct {
	Id             int
	TransactionId  *int
	UserId         int
	IdempotencyKey string
	Operation      string
	Status         string
	Response       []byte
	UpdatedAt      *time.Time
	CreatedAt      time.Time
}
