-- +goose Up
-- +goose StatementBegin
CREATE TABLE card_items
(
    id        UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    card_id   UUID REFERENCES cards (id) ON DELETE CASCADE NOT NULL,
    title_enc bytea,
    body_enc  bytea
);

CREATE INDEX ix_card_items_card_id
    ON card_items (card_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE card_items;
-- +goose StatementEnd
