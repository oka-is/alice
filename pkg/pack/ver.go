package pack

import (
	"github.com/oka-is/alice/lib/cryptos"
	"github.com/oka-is/srp6ago"
)

type VerByte byte

const (
	VerUnknown = VerByte(iota)
	Ver1
	Ver666
)

type Ver struct {
	MacByte          cryptos.MacByte
	AedCipherByte    cryptos.AedCipherByte
	PubCipherByte    cryptos.PubCipherByte
	DeriveByte       cryptos.DeriveByte
	SrpGroup         srp6ago.Params
	DeriveIterations int
}

var (
	ver1 = &Ver{
		MacByte:          cryptos.HmacSha256,
		AedCipherByte:    cryptos.AES256GCM,
		PubCipherByte:    cryptos.Rsa2048Sha256,
		DeriveByte:       cryptos.Pbkdf2Sha256,
		SrpGroup:         srp6ago.RFC5054b4096Sha256,
		DeriveIterations: 10_000,
	}

	// for testing purpose
	ver666 = &Ver{
		MacByte:          cryptos.HmacSha256,
		AedCipherByte:    cryptos.AES256GCM,
		PubCipherByte:    cryptos.Rsa1024Sha256,
		DeriveByte:       cryptos.Pbkdf2Sha256,
		SrpGroup:         srp6ago.RFC5054b1024Sha256,
		DeriveIterations: 1,
	}
)

func NewWer(byte VerByte) *Ver {
	switch byte {
	case Ver1:
		return ver1
	case Ver666:
		return ver666
	default:
		return nil
	}
}

func (v *Ver) NewSrpServer(verifier, salt []byte) *srp6ago.Server {
	return srp6ago.NewServer(verifier, salt, v.SrpGroup)
}
