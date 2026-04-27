CREATE TABLE IF NOT EXISTS USERS (
    id BIGINT generated always as identity primary key,
    firstname text not null,
    lastname text not null,
    birthday TIMESTAMPtz not null,
    phone_number text not null,
    password text not null,
    created_at TIMESTAMPtz not null default now()
);




