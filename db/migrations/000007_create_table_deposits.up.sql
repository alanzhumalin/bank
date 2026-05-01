DO $$
BEGIN
    CREATE type deposit_source as enum('terminal', 'bank');
EXCEPTION
    WHEN duplicate_object then null;
END $$;

DO $$
BEGIN
    CREATE TYPE deposit_status as enum ('pending', 'completed', 'failed');
EXCEPTION
    when duplicate_object then null;
END $$;

create table if not exists deposits(
    id BIGINT generated always as identity PRIMARY KEY,
    account_id BIGINT REFERENCES accounts(id),
    amount numeric(12,2) not null check (amount>0),
    source deposit_source not null,
    status deposit_status not null default 'pending',
    status_message text not null,
    created_at TIMESTAMPtz not null DEFAULT now()
);