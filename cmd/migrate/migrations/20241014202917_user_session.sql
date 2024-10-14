-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_account (
    id INTEGER NOT NULL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    email_verified INTEGER NOT NULL DEFAULT 0,
    totp_key BYTEA,
    recovery_code BYTEA NOT NULL
);

CREATE INDEX email_index ON user_account(email);

CREATE TABLE session (
    id TEXT NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES user_account(id),
    expires_at INTEGER NOT NULL,
    two_factor_verified INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE email_verification_request (
    id TEXT NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES user_account(id),
    email TEXT NOT NULL,
    code TEXT NOT NULL,
    expires_at INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_account;
DROP TABLE session;
DROP TABLE email_verification_request;
-- +goose StatementEnd
