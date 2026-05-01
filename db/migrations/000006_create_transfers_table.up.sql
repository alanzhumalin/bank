CREATE TABLE IF NOT EXISTS transfers(
    id BIGINT generated always as identity PRIMARY key,
    transaction_id BIGINT REFERENCES transactions(id),
    sender_account_id BIGINT not null REFERENCES accounts(id),
    receiver_account_id BIGINT not null REFERENCES accounts(id),
    currency_id BIGINT not null REFERENCES currencies(id),
    amount numeric(12,2) not null check (amount > 0)
);