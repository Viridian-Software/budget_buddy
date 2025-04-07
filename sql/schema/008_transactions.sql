-- +goose Up
ALTER TABLE transactions
ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
ADD COLUMN is_recurring BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE transactions
DROP COLUMN updated_at,
DROP COLUMN is_recurring;