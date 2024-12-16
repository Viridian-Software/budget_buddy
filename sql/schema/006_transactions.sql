-- +goose Up
CREATE TABLE transactions(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    account_id UUID NOT NULL,
    amount MONEY NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_account_id FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE transactions;