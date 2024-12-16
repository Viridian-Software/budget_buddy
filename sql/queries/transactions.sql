-- name: CreateTransaction :one
INSERT INTO transactions(
    id, created_at, user_id, account_id, amount
) VALUES (
    $1,
    NOW(),
    $2,
    $3,
    $4
) RETURNING *;