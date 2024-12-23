// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: accounts.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createNewAccount = `-- name: CreateNewAccount :one
INSERT INTO accounts(id, account_name, current_balance, account_type, user_id, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    NOW(),
    NOW()
) RETURNING id, account_name, current_balance, account_type, user_id, created_at, updated_at
`

type CreateNewAccountParams struct {
	AccountName    string
	CurrentBalance string
	AccountType    string
	UserID         uuid.UUID
}

func (q *Queries) CreateNewAccount(ctx context.Context, arg CreateNewAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createNewAccount,
		arg.AccountName,
		arg.CurrentBalance,
		arg.AccountType,
		arg.UserID,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.AccountName,
		&i.CurrentBalance,
		&i.AccountType,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, account_name, current_balance, account_type, user_id, created_at, updated_at FROM accounts
WHERE id = $1
`

func (q *Queries) GetAccount(ctx context.Context, id uuid.UUID) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.AccountName,
		&i.CurrentBalance,
		&i.AccountType,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAccountsByUser = `-- name: GetAccountsByUser :many
SELECT id, account_name, current_balance, account_type, user_id, created_at, updated_at FROM accounts
WHERE user_id = $1
ORDER BY created_at ASC
`

func (q *Queries) GetAccountsByUser(ctx context.Context, userID uuid.UUID) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, getAccountsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.AccountName,
			&i.CurrentBalance,
			&i.AccountType,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBalance = `-- name: UpdateBalance :one
UPDATE accounts
SET current_balance = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, account_name, current_balance, account_type, user_id, created_at, updated_at
`

type UpdateBalanceParams struct {
	ID             uuid.UUID
	CurrentBalance string
}

func (q *Queries) UpdateBalance(ctx context.Context, arg UpdateBalanceParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateBalance, arg.ID, arg.CurrentBalance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.AccountName,
		&i.CurrentBalance,
		&i.AccountType,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const verifyAccountExistence = `-- name: VerifyAccountExistence :exec
SELECT EXISTS (
    SELECT 1
    FROM accounts
    WHERE user_id = $1
)
`

func (q *Queries) VerifyAccountExistence(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, verifyAccountExistence, userID)
	return err
}
