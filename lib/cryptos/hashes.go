package cryptos

import (
	"crypto"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

const (
	HashUnknown = HashByte(iota)
	SHA256
	SHA512
)

func NewHash(byte HashByte) IHash {
	switch byte {
	case SHA256:
		return &Hash{b: SHA256, fn: HashCrypto(SHA256).New}
	case SHA512:
		return &Hash{b: SHA512, fn: HashCrypto(SHA512).New}
	default:
		return nil
	}
}

func HashSize(byte HashByte) int {
	switch byte {
	case SHA256:
		return sha256.Size
	case SHA512:
		return sha512.Size
	default:
		return 0
	}
}

func HashCrypto(byte HashByte) crypto.Hash {
	switch byte {
	case SHA256:
		return crypto.SHA256
	case SHA512:
		return crypto.SHA512
	default:
		return 0
	}
}

type Hash struct {
	b  HashByte
	fn func() hash.Hash
}

func (h *Hash) Byte() HashByte {
	return h.b
}

func (h *Hash) New() hash.Hash {
	return h.fn()
}

func (h *Hash) HashSum(input []byte) []byte {
	fn := h.New()
	fn.Write(input)
	return fn.Sum(nil)
}
