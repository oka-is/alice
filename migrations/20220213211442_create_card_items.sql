-- +goose Up
-- +goose StatementBegin
CREATE TABLE card_items
(
    id        UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    card_id   UUID REFERENCES cards (id) ON DELETE CASCADE NOT NULL,
    position  int                                          NOT NULL,
    title_enc bytea,
    body_enc  bytea,
    hidden    boolean
);

CREATE UNIQUE INDEX ux_card_items_card_id_position
    ON card_items (card_id, position);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE card_items;
-- +goose StatementEnd
