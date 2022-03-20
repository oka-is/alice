-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN otp_secret bytea;
ALTER TABLE users
    ADD COLUMN otp_candidate bytea;

ALTER TABLE sessions
    ADD COLUMN otp_succeed boolean;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN otp_secret;
ALTER TABLE users
    DROP COLUMN otp_candidate;

ALTER TABLE sessions
    DROP COLUMN otp_succeed;
-- +goose StatementEnd
