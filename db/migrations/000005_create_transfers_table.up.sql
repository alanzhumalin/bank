CREATE type transfer_status as enum ('pending', 'completed', 'failed');

CREATE TABLE transfers(
    id BIGINT generated always as identity PRIMARY key,
    sender_account_id BIGINT not null REFERENCES accounts(id),
    receiver_account_id BIGINT not null REFERENCES accounts(id),
    currency_id BIGINT not null REFERENCES currencies(id),
    amount numeric(12,2) not null check (amount > 0),
    status transfer_status not null DEFAULT 'pending',
    created_at TIMESTAMPtz not null default now()
);