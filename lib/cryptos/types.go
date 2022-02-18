package cryptos

import (
	"hash"
)

type (
	HashByte      byte
	MacByte       byte
	AedCipherByte byte
	PubCipherByte byte
	DeriveByte    byte
)

type (
	MacNewFN func(h func() hash.Hash, key []byte) hash.Hash
	DeriveFN func(password, salt []byte, iterations, keyLen int, h func() hash.Hash) []byte
)

type IHash interface {
	New() hash.Hash
	Byte() HashByte
	HashSum(input []byte) []byte
}

// IDerive interface for key derivation functions
type IDerive interface {
	Derive(password, salt []byte, iterations, len int) []byte
}

// IAedCipher interface for block cypher with authentication algorithms
type IAedCipher interface {
	Encrypt(plaintext, addon []byte) []byte
	Decrypt(ciphertext, addon []byte) ([]byte, error)
	Init(key, iv []byte) (IAedCipher, error)
}

type IPubCipher interface {
	GeneratePair() (priv []byte, pub []byte, err error)
	Encrypt(pub []byte, data []byte) ([]byte, error)
	Decrypt(priv []byte, data []byte) ([]byte, error)
	Sign(priv []byte, data []byte) ([]byte, error)
	Verify(pub []byte, data []byte, sig []byte) error
}

// IMac interface for Message Authentication Code
type IMac interface {
	Init(key []byte) IMac
	Write(data []byte) (int, error)
	Sum() []byte
	HashSum(input []byte) []byte
	Reset()
}
