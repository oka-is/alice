package jwt

import "errors"

var (
	ErrRequired     = errors.New("ERR_REQUIRED")
	ErrAlgoMismatch = errors.New("ERR_ALGO_MISMATCH")
	ErrInvalid      = errors.New("ERR_INVALID")
)
