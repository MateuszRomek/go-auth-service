-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_account (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    email_verified INTEGER NOT NULL DEFAULT 0,
    salt TEXT NOT NULL
);

CREATE INDEX email_index ON user_account(email);

CREATE TABLE session (
    id TEXT NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_account(id),
    expires_at BIGINT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE session;
DROP TABLE user_account;
-- +goose StatementEnd
