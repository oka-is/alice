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

func (s *SqlxDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return SQLErr(s.DB.GetContext(ctx, dest, query, args...))
}
