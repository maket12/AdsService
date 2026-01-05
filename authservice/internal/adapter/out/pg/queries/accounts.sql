-- name: CreateAccount :exec
INSERT INTO accounts (
    id,
    email,
    password_hash,
    status,
    email_verified,
    created_at,
    updated_at,
    last_login_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: GetAccountByEmail :one
SELECT
    id,
    email,
    password_hash,
    status,
    email_verified,
    created_at,
    updated_at,
    last_login_at
FROM accounts
WHERE email = $1;

-- name: GetAccountByID :one
SELECT
    id,
    email,
    password_hash,
    status,
    email_verified,
    created_at,
    updated_at,
    last_login_at
FROM accounts
WHERE id = $1;

-- name: UpdateAccountEmail :exec
UPDATE accounts
SET
    email = $2,
    updated_at = $3
WHERE id = $1;

-- name: UpdateAccountStatus :exec
UPDATE accounts
SET
    status = $2,
    updated_at = $3
WHERE id = $1;

-- name: UpdateAccountPassword :exec
UPDATE accounts
SET
    password_hash = $2,
    updated_at = $3
WHERE id = $1;

-- name: MarkAccountLogin :exec
UPDATE accounts
SET
    last_login_at = $2,
    updated_at = $3
WHERE id = $1;

-- name: VerifyAccountEmail :exec
UPDATE accounts
SET
    email_verified = true,
    updated_at = $2
WHERE id = $1;
