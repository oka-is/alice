package domain

type User struct {
	ID              NullString `db:"id"`
	Ver             NullInt64  `db:"ver"` // pack version
	Readonly        NullBool   `db:"readonly"`
	RestrictSharing NullBool   `db:"restrict_sharing"`
	Identity        NullString `db:"identity"` // hashed username
	Verifier        NullBytea  `db:"verifier"`
	SrpSalt         NullBytea  `db:"srp_salt"`
	PasswdSalt      NullBytea  `db:"passwd_salt"`
	PrivKeyEnc      NullBytea  `db:"priv_key_enc"`
	PubKey          NullBytea  `db:"pub_key"`
	OtpSecret       NullBytea  `db:"otp_secret"`
	OtpCandidate    NullBytea  `db:"otp_candidate"`
	CreatedAt       NullTime   `db:"created_at"`
}

func (u *User) IsOtpEnabled() bool {
	return len(u.OtpSecret.Bytea) > 0
}

func (u *User) IsOtpDisabled() bool {
	return !u.IsOtpEnabled()
}
