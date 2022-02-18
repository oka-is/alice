-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions
(
    jti          UUID PRIMARY KEY,
    user_id      UUID REFERENCES users ON DELETE CASCADE,
    candidate_id UUID REFERENCES users ON DELETE CASCADE,
    srp_state    bytea,
    time_from    TIMESTAMPTZ NOT NULL,
    time_to      TIMESTAMPTZ NOT NULL
);
CREATE INDEX ix_sessions_time ON sessions (time_to);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
