package cryptos

import (
	"crypto/hmac"
	"hash"
)

const (
	MacUnknown = MacByte(iota)
	HmacSha256
)

type Mac struct {
	macFN MacNewFN
	actor hash.Hash
	hash  IHash
}

func (m *Mac) Init(key []byte) IMac {
	m.actor = hmac.New(m.hash.New, key)
	return m
}

func (m *Mac) Write(data []byte) (int, error) {
	return m.actor.Write(data)
}

func (m *Mac) Sum() []byte {
	return m.actor.Sum(nil)
}

func (m *Mac) Reset() {
	m.actor.Reset()
}

func (m *Mac) HashSum(input []byte) []byte {
	m.Reset()
	_, _ = m.actor.Write(input)
	return m.Sum()
}

func NewMac(byte MacByte) IMac {
	switch byte {
	case HmacSha256:
		return &Mac{macFN: hmac.New, hash: NewHash(SHA256)}
	default:
		return nil
	}
}

func MacSize(byte MacByte) int {
	switch byte {
	case HmacSha256:
		return HashSize(SHA256)
	default:
		return -0
	}
}
