-- +goose Up
-- +goose StatementBegin
CREATE TABLE cards
(
    id           UUID                                                       DEFAULT gen_random_uuid() PRIMARY KEY,
    workspace_id UUID REFERENCES workspaces (id) ON DELETE CASCADE NOT NULL,
    archived     boolean,
    title_enc    bytea,
    tags_enc     bytea[],
    created_at   TIMESTAMPTZ                                       NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cards;
-- +goose StatementEnd
