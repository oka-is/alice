-- +goose Up
-- +goose StatementBegin
ALTER TABLE sessions
    ADD COLUMN otp_attempts smallint;
CREATE INDEX ix_sessions_candidate_time_from
    ON sessions (candidate_id, time_from) WHERE candidate_id IS NOT NULL;
CREATE INDEX ix_sessions_user_time_from
    ON sessions (user_id, time_from) WHERE user_id IS NOT NULL;
ALTER INDEX ix_sessions_time RENAME TO ix_sessions_time_to;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sessions
    DROP COLUMN otp_attempts;
ALTER INDEX ix_sessions_time_to RENAME TO ix_sessions_time;
DROP INDEX ix_sessions_candidate_time_from;
DROP INDEX ix_sessions_user_time_from;
-- +goose StatementEnd
