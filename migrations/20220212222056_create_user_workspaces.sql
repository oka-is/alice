-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_workspaces
(
    id           UUID                                                       DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id      UUID REFERENCES users (id) ON DELETE CASCADE      NOT NULL,
    owner_id     UUID REFERENCES users (id) ON DELETE CASCADE      NOT NULL,
    workspace_id UUID REFERENCES workspaces (id) ON DELETE CASCADE NOT NULL,
    aed_key_enc  bytea                                             NOT NULL,
    aed_key_tag  bytea                                             NOT NULL,
    created_at   TIMESTAMPTZ                                       NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ux_user_workspaces_user_id
    ON user_workspaces (user_id, workspace_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_workspaces;
-- +goose StatementEnd
