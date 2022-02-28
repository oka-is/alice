package storage

import (
	"context"
	"fmt"
	"strings"
)

// Truncate empty a table or set of tables
func (s *Storage) Truncate(ctx context.Context, tables ...string) error {
	statements := make([]string, len(tables))

	for ix, table := range tables {
		statements[ix] = fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table)
	}

	_, err := s.db.Exec(strings.Join(statements, ";\n"))
	return err
}

// TruncateAll empty whole database
func (s *Storage) TruncateAll(ctx context.Context) error {
	tables := make([]string, 0)
	query := Builder().
		Select("table_name").
		From("information_schema.tables").
		Where("table_schema = 'public'").
		Where("table_name NOT IN ('goose_db_version')").
		Where("table_type = 'BASE TABLE'")

	err := s.Select(ctx, s.db, &tables, query)
	if err != nil {
		return err
	}

	return s.Truncate(ctx, tables...)
}
