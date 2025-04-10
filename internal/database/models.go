// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID             uuid.UUID
	AccountName    string
	CurrentBalance string
	AccountType    string
	UserID         uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RefreshToken struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt sql.NullTime
}

type Transaction struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UserID      uuid.UUID
	AccountID   uuid.UUID
	Amount      string
	Description sql.NullString
	UpdatedAt   time.Time
	IsRecurring bool
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Email          string
	IsAdmin        bool
	FirstName      string
	LastName       string
	HashedPassword string
}
