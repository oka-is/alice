package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/oka-is/alice/lib/jwt"
	"github.com/oka-is/alice/pkg/domain"
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
	query := Builder().
		Update("sessions").
		Set("candidate_id", candidateID).
		Set("srp_state", srp).
		Where("jti = ?", jti)

	_, err := s.Exec(ctx, s.db, query)
	return err
}

func (s *Storage) NominateSession(ctx context.Context, jti string) error {
	query := Builder().
		Update("sessions").
		Set("user_id", Expr("candidate_id")).
		Set("srp_state", Expr("NULL")).
		Where("jti = ?", jti)

	_, err := s.Exec(ctx, s.db, query)
	return err
}

func (s *Storage) DeleteSession(ctx context.Context, jti string) error {
	query := Builder().Delete("sessions").Where("jti = ?", jti)
	return s.Exec1(ctx, s.db, query)
}

// FindSession test case usage only, for regular usage, please use RetrieveSession
func (s *Storage) FindSession(ctx context.Context, jti string) (out domain.Session, err error) {
	query := Builder().Select("*").From("sessions").Where("jti = ?", jti)
	err = s.Get(ctx, s.db, &out, query)
	return
}

func (s *Storage) insertSession(ctx context.Context, db IConn, session *domain.Session) error {
	query := Builder().
		Insert("sessions").
		Columns("jti", "time_from", "time_to").
		Values(session.Jti, session.TimeFrom, session.TimeTo)

	return s.Exec1(ctx, db, query)
}
