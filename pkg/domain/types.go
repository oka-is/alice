package domain

import (
	"database/sql"
	"time"

	"github.com/wault-pw/alice/lib/null"
)

type (
	NullString     = sql.NullString
	NullBool       = sql.NullBool
	NullTime       = sql.NullTime
	NullInt64      = sql.NullInt64
	NullBytea      = null.Bytea
	NullByteaSlice = null.ByteaSlice
)

func NewEmptyString(input string) NullString {
	return NullString{String: input, Valid: input != ""}
}

func NewNullString() NullString {
	return NullString{Valid: false}
}

func NewNullBytea() NullBytea {
	return NullBytea{Bytea: []byte{}, Valid: false}
}

func NewEmptyInt64(input int64) NullInt64 {
	return NullInt64{Int64: input, Valid: input != 0}
}

func NewNullInt64(input int64) NullInt64 {
	return NullInt64{Int64: input, Valid: true}
}

func NewEmptyBytes(input []byte) NullBytea {
	return NullBytea{Bytea: input, Valid: len(input) != 0}
}

func NewEmptyByteSlice(input [][]byte) NullByteaSlice {
	return NullByteaSlice{Slice: input, Valid: len(input) != 0}
}

func NewEmptyTime(input time.Time) NullTime {
	return NullTime{Time: input, Valid: !input.IsZero()}
}

func NewEmptyBool(input bool) NullBool {
	return NullBool{Bool: input, Valid: true}
}
