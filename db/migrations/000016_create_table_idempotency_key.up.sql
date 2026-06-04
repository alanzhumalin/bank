DO $$
BEGIN
    create type operation_type as enum ('withdraw', 'deposit', 'transfer');
EXCEPTION
    when duplicate_object then NULL;
END $$;

create table if not exists idempotency_keys(
    id BIGINT generated always as identity PRIMARY KEY,
    transaction_id BIGINT REFERENCES transactions(id),
    user_id BIGINT REFERENCES users(id) not null,
    idempotency_key text not null,
    operation operation_type not null,
    status text not null DEFAULT 'pending',
    response JSONB,
    created_at TIMESTAMPTZ not null DEFAULT now(),
    updated_at TIMESTAMPTZ,
    UNIQUE(user_id, idempotency_key)
);