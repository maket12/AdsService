DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_type') THEN
            CREATE TYPE role_type AS ENUM ('user', 'admin');
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS account_roles (
    account_id uuid NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    role role_type NOT NULL,
    PRIMARY KEY (account_id)
);

CREATE INDEX idx_account_roles_role ON account_roles(role);