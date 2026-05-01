DO $$
BEGIN
    create type transaction_type as enum ('transfer', 'deposit','withdraw');
EXCEPTION
    when duplicate_object then null;
END $$;


DO $$
BEGIN
    create type transaction_status as enum('pending', 'failed', 'completed');
EXCEPTION
    when duplicate_object then null;
END $$;


create table if not exists transactions(
    id BIGINT generated always as identity PRIMARY KEY,
    type transaction_type not null,
    amount numeric(12,2) not null check (amount >0),
    account_id BIGINT not null REFERENCES accounts(id),
    status transaction_status not null DEFAULT 'pending',
    status_message text not null,
    created_at TIMESTAMPtz not null DEFAULT now()
);