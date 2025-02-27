-- +goose Up
ALTER TABLE transactions
ADD COLUMN description TEXT;

-- +goose Down
ALTER TABLE transactions
DROP COLUMN description;