-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id           UUID                 DEFAULT gen_random_uuid() PRIMARY KEY,
    ver          INT         NOT NULL,
    identity     bytea       NOT NULL,
    verifier     bytea       NOT NULL,
    srp_salt     bytea       NOT NULL,
    passwd_salt  bytea       NOT NULL,
    priv_key_enc bytea       NOT NULL,
    pub_key      bytea       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ux_users_identity ON users (identity);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
