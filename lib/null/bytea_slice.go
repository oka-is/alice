package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
)

const (
	pgHexEsc = `\\x`
)

type ByteaSlice struct {
	Slice [][]byte
	Valid bool
}

// Scan implements the Scanner interface.
func (b *ByteaSlice) Scan(value interface{}) error {
	if value == nil {
		b.Slice, b.Valid = make([][]byte, 0), false
		return nil
	}

	got, err := parsePgArrayHex(value.(string))
	if err != nil {
		return fmt.Errorf("failed to decode pg array: %w", err)
	}

	b.Valid = true
	b.Slice = got

	return nil
}

// Value implements the driver Valuer interface.
func (b ByteaSlice) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Slice, nil
}

// `{"\\x3a86789526fc792b93b5c01ce78425e317eccfc6fdd51f5f9ce790d4136fc263252e"}`
func parsePgArrayHex(input string) ([][]byte, error) {
	buff := new(bytes.Buffer)
	out := make([][]byte, 0)

	for _, r := range input {
		switch r {
		case '{', '"':
			continue
		case '}', ',':
			got, err := hex.DecodeString(buff.String()[len(pgHexEsc):])
			if err != nil {
				return nil, err
			}
			out = append(out, got)
			buff.Reset()
		default:
			buff.WriteRune(r)
		}
	}

	return out, nil
}
