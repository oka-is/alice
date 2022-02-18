package cryptos

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

const (
	AedCipherUnknown = AedCipherByte(iota)
	AES256GCM
)

func NewAedCipher(byte AedCipherByte) IAedCipher {
	switch byte {
	case AES256GCM:
		ks, is, _ := AedCipherSizes(byte)
		return NewAesAedCipher(ks, is)
	default:
		return nil
	}
}

type AesAedCipher struct {
	keySize, ivSize int
	iv              []byte
	aed             cipher.AEAD
}

func NewAesAedCipher(keySize, ivSize int) IAedCipher {
	return &AesAedCipher{
		keySize: keySize,
		ivSize:  ivSize,
	}
}

func (a *AesAedCipher) Init(key, iv []byte) (IAedCipher, error) {
	block, err := aes.NewCipher(key)

	switch {
	case err != nil:
		return a, fmt.Errorf("failed to init block: %w", err)
	case len(key) != a.keySize:
		return a, fmt.Errorf("%w: broken key size, want: %d, got: %d", ErrBrokenSize, a.keySize, len(key))
	case len(iv) != a.ivSize:
		return a, fmt.Errorf("%w: broken iv size, want: %d, got: %d", ErrBrokenSize, a.ivSize, len(iv))
	}

	aed, err := cipher.NewGCM(block)
	if err != nil {
		return a, fmt.Errorf("failed to init AED: %w", err)
	}

	a.iv = iv
	a.aed = aed

	return a, nil
}

// Encrypt encrypts and signs data
func (a *AesAedCipher) Encrypt(input, addon []byte) []byte {
	return a.aed.Seal(nil, a.iv, input, addon)
}

// Decrypt decrypts & verifies of data
func (a *AesAedCipher) Decrypt(input, addon []byte) ([]byte, error) {
	return a.aed.Open(nil, a.iv, input, addon)
}

func AedCipherSizes(byte AedCipherByte) (key, iv, block int) {
	switch byte {
	case AES256GCM:
		return 256 >> 3, 12, aes.BlockSize
	default:
		return -1, -1, -1
	}
}
