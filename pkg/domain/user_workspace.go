package domain

type UserWorkspace struct {
	ID          NullString `db:"id"`
	UserID      NullString `db:"user_id"`
	OwnerID     NullString `db:"owner_id"`
	WorkspaceID NullString `db:"workspace_id"`
	AedKeyEnc   NullBytea  `db:"scrt_key_enc"`
	CreatedAt   NullTime   `db:"created_at"`
}
