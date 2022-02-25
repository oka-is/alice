package storage

import (
	"context"
)

func (s *Storage) Ping(ctx context.Context) error {
	query := Builder().Select("1")
	_, err := s.Exec(ctx, s.db, query)
	return err
}
