-- +goose Up
ALTER TABLE users
ADD first_name TEXT NOT NULL;
ALTER TABLE users
ADD last_name TEXT NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN first_name, last_name;