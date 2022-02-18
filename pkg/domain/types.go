package domain

import (
	"database/sql"
	"time"

	"github.com/oka-is/alice/lib/null"
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

func NewEmptyInt64(input int64) NullInt64 {
	return NullInt64{Int64: input, Valid: input != 0}
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
