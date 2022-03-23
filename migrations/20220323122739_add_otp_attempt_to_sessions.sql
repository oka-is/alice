-- +goose Up
-- +goose StatementBegin
ALTER TABLE sessions ADD COLUMN otp_attempts smallint NOT NULL DEFAULT 0;
CREATE INDEX ix_sessions_candidate_time_from ON sessions(candidate_id, time_from);
ALTER INDEX ix_sessions_time RENAME TO ix_sessions_time_to;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sessions DROP COLUMN otp_attempts;
ALTER INDEX ix_sessions_time_to RENAME TO ix_sessions_time;
-- +goose StatementEnd
