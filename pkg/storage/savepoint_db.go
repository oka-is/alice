package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// SavepointDB wrapper around sqlx database
type SavepointDB struct {
	db  *sqlx.DB
	uid string
	tx  ITransaction
	mx  sync.Mutex
}

func NewSavepointDB(db *sqlx.DB) *SavepointDB {
	return &SavepointDB{db: db}
}

func (s *SavepointDB) SqlDB() *sql.DB {
	return s.db.DB
}

func (s *SavepointDB) BeginTxx(ctx context.Context, opts *TxOpts) (ITransaction, error) {
	uid := strings.Replace(uuid.New().String(), "-", "", -1)
	uid = fmt.Sprintf("test_%s", uid)
	point := &SavepointDB{db: s.db, uid: uid}

	if _, err := s.tx.ExecContext(context.Background(), fmt.Sprintf("SAVEPOINT %s", uid)); err != nil {
		return point, fmt.Errorf("failed to create a savepoint: %w", err)
	}

	return point, nil
}

func (s *SavepointDB) ensure() {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.tx == nil {
		s.tx = NewSavepointTX(s.db)
	}
}

// Flush will roll back parent transaction
// (all save points inside will be destroyed)
func (s *SavepointDB) Flush() {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.tx == nil {
		return
	}

	if err := s.tx.Rollback(); err != nil {
		panic(fmt.Errorf("failed to tollback a savepoint transaction: %w", err))
	}

	s.tx = nil
}

func (s *SavepointDB) Rollback() error {
	s.ensure()

	// ensure we rollbacks hole transaction on SQL error,
	// to make test database clean
	_, err := s.tx.ExecContext(context.Background(), fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", s.uid))
	if err != nil {
		return err
	}

	s.Flush()
	return nil
}

func (s *SavepointDB) Commit() error {
	s.ensure()
	_, err := s.tx.ExecContext(context.Background(), fmt.Sprintf("RELEASE SAVEPOINT %s", s.uid))
	return err
}

func (s *SavepointDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	s.ensure()
	return s.tx.ExecContext(ctx, query, args...)
}

func (s *SavepointDB) SelectContext(ctx context.Context, des interface{}, query string, args ...interface{}) error {
	s.ensure()
	return s.tx.SelectContext(ctx, des, query, args...)
}

func (s *SavepointDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	s.ensure()
	return SQLErr(s.tx.GetContext(ctx, dest, query, args...))
}

func (s *SavepointDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	s.ensure()
	return s.tx.QueryRowContext(ctx, query, args...)
}

func NewSavepointTX(db *sqlx.DB) *sqlx.Tx {
	tx, err := db.BeginTxx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(fmt.Errorf("filed to begin a transaction for a savepoint: %w", err))
	}
	return tx
}
