package domain

type CardItem struct {
	ID       NullString `db:"id"`
	CardID   NullString `db:"card_id"`
	Position NullInt64  `db:"position"`
	TitleEnc NullBytea  `db:"title_enc"`
	BodyEnc  NullBytea  `db:"body_enc"`
	Hidden   NullBool   `db:"hidden"`
}
