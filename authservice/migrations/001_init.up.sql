CREATE EXTENSION IF NOT EXISTS citext;

DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'account_status') THEN
            CREATE TYPE account_status AS ENUM ('active', 'blocked', 'deleted');
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS accounts (
    id uuid PRIMARY KEY,
    email citext NOT NULL UNIQUE,
    password_hash text NOT NULL,
    status account_status NOT NULL DEFAULT 'active',
    email_verified boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    last_login_at timestamptz
);

CREATE INDEX idx_accounts_status ON accounts(status);