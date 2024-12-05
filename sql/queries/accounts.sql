-- name: CreateNewAccount :one
INSERT INTO accounts(id, account_name, current_balance, account_type, user_id, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    NOW(),
    NOW()
) RETURNING *;

-- name: GetAccountsByUser :many
SELECT * FROM accounts
WHERE user_id = $1
ORDER BY created_at ASC;