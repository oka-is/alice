package storage

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/oka-is/alice/pkg/validator"
)

type Storage struct {
	db        *sqlx.DB
	validator validator.IValidator
}

func Connect(dsn string, validator validator.IValidator) (*Storage, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("filed to connect to postgres: %w", err)
	}

	return &Storage{db: db, validator: validator}, nil
}

func (s *Storage) SetValidator(validator validator.IValidator) *Storage {
	s.validator = validator
	return s
}

func (s *Storage) SqlDB() *sql.DB {
	return s.db.DB
}

// Builder returns SQL Builder object
func Builder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
