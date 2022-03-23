package domain

type Session struct {
	Jti         NullString `db:"jti"`     // JWT token ID
	UserID      NullString `db:"user_id"` // might be NULL if the user is not authenticated yet
	CandidateID NullString `db:"candidate_id"`
	SrpState    NullBytea  `db:"srp_state"`
	OtpSucceed  NullBool   `db:"otp_succeed"`
	OtpAttempts NullInt64  `db:"otp_attempts"`
	TimeFrom    NullTime   `db:"time_from"`
	TimeTo      NullTime   `db:"time_to"`
}
