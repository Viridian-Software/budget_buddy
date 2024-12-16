// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: transactions.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions(
    id, created_at, user_id, account_id, amount
) VALUES (
    $1,
    NOW(),
    $2,
    $3,
    $4
) RETURNING id, created_at, user_id, account_id, amount
`

type CreateTransactionParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	AccountID uuid.UUID
	Amount    string
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.ID,
		arg.UserID,
		arg.AccountID,
		arg.Amount,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UserID,
		&i.AccountID,
		&i.Amount,
	)
	return i, err
}