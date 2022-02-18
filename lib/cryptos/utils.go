package cryptos

import (
	"crypto/rand"
	"crypto/subtle"
)

func SecureRand(buff []byte) error {
	_, err := rand.Read(buff)
	return err
}

func SecureRandSize(size int) ([]byte, error) {
	buff := make([]byte, size)
	err := SecureRand(buff)
	return buff, err
}

func ConstantTimeByteEq(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
