create table if not exists sessions(
    id uuid primary key not null,
    hashed_refresh_token text not null unique,
    user_id bigint not null references users(id) on delete CASCADE, 
    device text not null,
    ip inet not null,
    created_at timestamptz not null DEFAULT now(),
    expires_at timestamptz not null,
    is_active boolean not null DEFAULT true
);