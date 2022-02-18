-- +goose Up
-- +goose StatementBegin
CREATE TABLE workspaces
(
    id         UUID                 DEFAULT gen_random_uuid() PRIMARY KEY,
    title_enc  bytea,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workspaces;
-- +goose StatementEnd
