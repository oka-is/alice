package domain

type UserWithWorkspace struct {
	RecordID           NullString `db:"record_id"`
	UserID             NullString `db:"user_id"`
	OwnerID            NullString `db:"owner_id"`
	OwnerPubKey        NullBytea  `db:"owner_pub_key"`
	WorkspaceID        NullString `db:"workspace_id"`
	AedKeyEnc          NullBytea  `db:"aed_key_enc"`
	TitleEnc           NullBytea  `db:"title_enc"`
	RecordCreatedAt    NullTime   `db:"record_created_at"`
	WorkspaceCreatedAt NullTime   `db:"workspace_created_at"`
}
