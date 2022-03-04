package storage

import (
	"database/sql"
	"errors"
)

var (
	// ErrNotFound is an error for cases when requested entity is not found.
	ErrNotFound = errors.New("entity not found")
)

// SQLErr wraps error from database
func SQLErr(err error) error {
	switch err {
	case sql.ErrNoRows:
		return ErrNotFound
	default:
		return err
	}
}
