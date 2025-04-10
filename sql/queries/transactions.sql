-- name: CreateTransaction :one
INSERT INTO transactions(
    id, created_at, user_id, account_id, amount, description, updated_at, is_recurring
) VALUES (
    gen_random_uuid(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    NOW(),
    $5
) RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions WHERE id = $1;

-- name: GetAllTransactions :many
SELECT * FROM transactions
WHERE account_id = $1;

-- name: GetTransactionByID :one
SELECT * FROM transactions
WHERE id = $1;