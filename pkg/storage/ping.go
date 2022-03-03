package storage

import (
	"context"
)

func (s *Storage) Ping(ctx context.Context) error {
	return s.Exec1(ctx, s.db, Builder().Select("1"))
}
