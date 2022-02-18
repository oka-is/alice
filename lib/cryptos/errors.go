package cryptos

import "errors"

var (
	ErrBrokenSize      = errors.New("ERR_BROKEN_SIZE")
	ErrBrokenSignature = errors.New("ERR_BROKEN_SIGNATURE")
)
