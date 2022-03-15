package domain

type User struct {
	ID         NullString `db:"id"`
	Ver        NullInt64  `db:"ver"` // pack version
	Readonly   NullBool   `db:"readonly"`
	Identity   NullString `db:"identity"` // hashed username
	Verifier   NullBytea  `db:"verifier"`
	SrpSalt    NullBytea  `db:"srp_salt"`
	PasswdSalt NullBytea  `db:"passwd_salt"`
	PrivKeyEnc NullBytea  `db:"priv_key_enc"`
	PubKey     NullBytea  `db:"pub_key"`
	CreatedAt  NullTime   `db:"created_at"`
}
