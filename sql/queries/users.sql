-- name: GetAllUsers :many
SELECT * FROM users;

-- name: AddUser :one
INSERT INTO users(id, created_at, updated_at, email, is_admin, first_name, last_name, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    FALSE,
    $2,
    $3,
    $4
) RETURNING id, created_at, updated_at, email, is_admin, first_name, last_name;

-- name: DeleteAllUsers :exec
DELETE FROM users;