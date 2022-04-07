-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN restrict_sharing boolean;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN restrict_sharing;
-- +goose StatementEnd
