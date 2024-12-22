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

-- name: GetAllUsers :many
SELECT * FROM users
ORDER BY created_at ASC;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: SetUserAsAdmin :one
UPDATE users
SET is_admin = TRUE
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_admin, first_name, last_name;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUserInformation :one
UPDATE users
SET updated_at = NOW(), email = $1, first_name = $2, last_name = $3
WHERE id = $4
RETURNING id, created_at, updated_at, email, first_name, last_name;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;