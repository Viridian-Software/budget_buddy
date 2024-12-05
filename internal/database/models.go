// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
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
