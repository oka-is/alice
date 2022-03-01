package domain

type Card struct {
	ID          NullString     `db:"id"`
	WorkspaceID NullString     `db:"workspace_id"`
	Archived    NullBool       `db:"archived"`
	TitleEnc    NullBytea      `db:"title_enc"`
	TagsEnc     NullByteaSlice `db:"tags_enc"`
	CreatedAt   NullTime       `db:"created_at"`
}
