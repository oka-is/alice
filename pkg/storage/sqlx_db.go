package storage

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// SqlxDB wrapper around sqlx database
type SqlxDB struct {
	*sqlx.DB
}

func NewSqlxDB(db *sqlx.DB) IDb {
	return &SqlxDB{db}
}

func (s *SqlxDB) SqlDB() *sql.DB {
	return s.DB.DB
}

func (s *SqlxDB) BeginTxx(ctx context.Context, opts *TxOpts) (ITransaction, error) {
	return s.DB.BeginTxx(ctx, opts)
}
