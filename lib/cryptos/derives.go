package cryptos

import (
	"golang.org/x/crypto/pbkdf2"
)

const (
	DeriveUnknown = DeriveByte(iota)
	Pbkdf2Sha256
)

func NewDerive(byte DeriveByte) IDerive {
	switch byte {
	case Pbkdf2Sha256:
		return &Derive{hash: NewHash(SHA256), fn: pbkdf2.Key}
	default:
		return nil
	}
}

type Derive struct {
	hash IHash
	fn   DeriveFN
}

// Derive implements IDerive
func (d *Derive) Derive(password, salt []byte, iterations, len int) []byte {
	return d.fn(password, salt, iterations, len, d.hash.New)
}

func DeriveSaltSize(byte DeriveByte) int {
	switch byte {
	case Pbkdf2Sha256:
		return HashSize(SHA256)
	default:
		return -0
	}
}
