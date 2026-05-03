do $$
BEGIN
    create type withdraw_source as enum ('terminal', 'bank');
EXCEPTION
    when duplicate_object then null;
End $$;

create table if not exists withdrawals(
    id BIGINT generated always as identity PRIMARY KEY,
    transaction_id BIGINT REFERENCES transactions(id) not null,
    account_id BIGINT REFERENCES accounts(id) not NULL,
    amount numeric(12,2) not null check (amount > 0),
    source withdraw_source not null
);