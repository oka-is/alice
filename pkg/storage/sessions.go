package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/wault-pw/alice/lib/jwt"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/pkg/validator"
)

const (
	SessionExpirationDuration = 5 * time.Hour
)

// IssueSession creates a session for anonymous and generates JWT token
func (s *Storage) IssueSession(ctx context.Context, opts jwt.IOts) (domain.Session, string, error) {
	session := domain.Session{
		Jti:      domain.NewEmptyString(domain.NewUUID()),
		TimeFrom: domain.NewEmptyTime(time.Now()),
		TimeTo:   domain.NewEmptyTime(time.Now().Add(SessionExpirationDuration)),
	}

	opts = opts.SetJti(session.Jti.String).SetExp(session.TimeTo.Time)
	token, err := opts.Marshall()
	if err != nil {
		return session, token, fmt.Errorf("jwt failed: %w", err)
	}

	return session, token, s.insertSession(ctx, s.db, &session)
}

// RetrieveSession find & verify a session by JWT token
func (s *Storage) RetrieveSession(ctx context.Context, opts jwt.IOts, token string) (session domain.Session, err error) {
	jti, err := opts.Unmarshall(token)
	if err != nil {
		return session, fmt.Errorf("jwt failed: %w", err)
	}

	query := Builder().Select("*").From("sessions").Where("jti = ?", jti).Limit(1)
	err = s.Get(ctx, s.db, &session, query)
	return
}

func (s *Storage) CandidateSession(ctx context.Context, jti, candidateID string, srp []byte) error {
	return s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		return s.candidateSession(c, tx, jti, candidateID, srp)
	})
}

// NominateSession moves candidateID to userID, from now it means that
// user is authorized
func (s *Storage) NominateSession(ctx context.Context, jti string) error {
	query := Builder().
		Update("sessions").
		Set("user_id", Expr("candidate_id")).
		Set("candidate_id", Expr("NULL")).
		Set("srp_state", Expr("NULL")).
		Where("jti = ?", jti)

	return s.Exec1(ctx, s.db, query)
}

func (s *Storage) OtpSessionSucceed(ctx context.Context, jti string) error {
	query := Builder().Update("sessions").Set("otp_succeed", true).Where("jti = ?", jti)
	return s.Exec1(ctx, s.db, query)
}

func (s *Storage) DeleteSession(ctx context.Context, jti string) error {
	query := Builder().Delete("sessions").Where("jti = ?", jti)
	return s.Exec1(ctx, s.db, query)
}

// DeleteUserSessionExcept deletes all issued sessions for specific user except the current one
// useful in password change or OTP enable
func (s *Storage) DeleteUserSessionExcept(ctx context.Context, userID, jti string) error {
	query := Builder().Delete("sessions").Where("user_id = ?", userID).Where("jti <> ?", jti)
	return s.Exec1(ctx, s.db, query)
}

// FindSession test case usage only, for regular usage, please use RetrieveSession
func (s *Storage) FindSession(ctx context.Context, jti string) (out domain.Session, err error) {
	query := Builder().Select("*").From("sessions").Where("jti = ?", jti)
	err = s.Get(ctx, s.db, &out, query)
	return
}

// MakeOtpAttempt do several thinks:
// 1) validates id user does not reach max OTP attempts per minute
// 1) increments OTP attempts counter
func (s *Storage) MakeOtpAttempt(ctx context.Context, jti string) error {
	return s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		return s.makeOtpAttempt(c, tx, jti)
	})
}

func (s *Storage) insertSession(ctx context.Context, db IConn, session *domain.Session) error {
	query := Builder().
		Insert("sessions").
		Columns(
			"jti",
			"user_id",
			"candidate_id",
			"otp_attempts",
			"time_from",
			"time_to").
		Values(
			session.Jti,
			session.UserID,
			session.CandidateID,
			session.OtpAttempts,
			session.TimeFrom,
			session.TimeTo)

	return s.Exec1(ctx, db, query)
}

func (s *Storage) candidateSession(ctx context.Context, db IConn, jti, candidateID string, srp []byte) error {
	loginAttempts, err := s.countLoginAttempts(ctx, db, candidateID, validator.LoginAttemptsDur)
	if err != nil {
		return fmt.Errorf("failed to count login attempts: %w", err)
	}

	err = s.validator.ValidateCandidateSession(validator.ValidateCandidateSessionOpts{
		Attempts: loginAttempts,
	})
	if err != nil {
		return fmt.Errorf("failed to validate candidate session: %w", err)
	}

	query := Builder().
		Update("sessions").
		Set("candidate_id", candidateID).
		Set("srp_state", srp).
		Where("jti = ?", jti)

	return s.Exec1(ctx, db, query)
}

// countLoginAttempts counts unsuccessful login attempts for the passed duration
func (s *Storage) countLoginAttempts(ctx context.Context, db IConn, candidateID string, dur time.Duration) (counter int, err error) {
	query := Builder().
		Select("COUNT(*)").
		From("sessions").
		Where("candidate_id = ?", candidateID).
		Where("time_from >= ?", time.Now().Add(-dur))

	err = s.Get(ctx, db, &counter, query)
	return
}

func (s *Storage) makeOtpAttempt(ctx context.Context, db IConn, jti string) error {
	counter, err := s.countOtpAttempts(ctx, db, jti, validator.OtpAttemptsDur)
	if err != nil {
		return fmt.Errorf("faild to count otp attempts: %w", err)
	}

	err = s.validator.ValidateOtpAttempt(validator.ValidateOtpAttemptOpts{
		Attempts: counter,
	})
	if err != nil {
		return fmt.Errorf("failed to validate otp attempt: %w", err)
	}

	return s.incrementOtpAttempts(ctx, db, jti)
}

func (s *Storage) countOtpAttempts(ctx context.Context, db IConn, jti string, dur time.Duration) (counter int, err error) {
	query := Builder().
		Select("COALESCE(SUM(otp_attempts), 0)").
		From("sessions").
		Where("user_id = (SELECT user_id FROM sessions WHERE jti = ? LIMIT 1)", jti).
		Where("time_from >= ?", time.Now().Add(-dur))

	err = s.Get(ctx, db, &counter, query)
	return
}

func (s *Storage) incrementOtpAttempts(ctx context.Context, db IConn, jti string) error {
	query := Builder().
		Update("sessions").
		Set("otp_attempts", Expr("COALESCE(otp_attempts, 0) + 1")).
		Where("jti = ?", jti)

	return s.Exec1(ctx, db, query)
}
