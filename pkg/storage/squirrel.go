package storage

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type (
	Sqlizer       = sq.Sqlizer
	SelectBuilder = sq.SelectBuilder
)

func Expr(sql string, args ...interface{}) Sqlizer {
	return sq.Expr(sql, args...)
}

func (s *Storage) Exec(ctx context.Context, conn IConn, query IBuilder) (Result, error) {
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("builder failed: %w", err)
	}
	return conn.ExecContext(ctx, sql, args...)
}

func (s *Storage) Exec1(ctx context.Context, conn IConn, query IBuilder) error {
	_, err := s.Exec(ctx, conn, query)
	return err
}

func (s *Storage) Select(ctx context.Context, conn IConn, dest interface{}, query IBuilder) error {
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("builder failed: %w", err)
	}
	return conn.SelectContext(ctx, dest, sql, args...)
}

func (s *Storage) Get(ctx context.Context, conn IConn, dest interface{}, query IBuilder) error {
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("builder failed: %w", err)
	}
	return conn.GetContext(ctx, dest, sql, args...)
}

func (s *Storage) QueryRow(ctx context.Context, conn IConn, query IBuilder) *Row {
	sql, args, err := query.ToSql()
	if err != nil {
		return &Row{}
	}
	return conn.QueryRowContext(ctx, sql, args...)
}

func (s *Storage) Tx(ctx context.Context, opts *TxOpts, fn TxFunc) (err error) {
	tx, err := s.db.BeginTxx(ctx, opts)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				err = fmt.Errorf("cannot commit transaction: %w", cmErr)
			}
		}
	}()

	err = fn(ctx, tx)
	return
}
