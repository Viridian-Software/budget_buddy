-- +goose Up
CREATE TABLE accounts(
    id UUID PRIMARY KEY,
    account_name TEXT NOT NULL,
    current_balance MONEY NOT NULL,
    account_type TEXT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE accounts;