package domain

type Workspace struct {
	ID        NullString `db:"id"`
	TitleEnc  NullBytea  `db:"title_enc"`
	CreatedAt NullTime   `db:"created_at"`
}
