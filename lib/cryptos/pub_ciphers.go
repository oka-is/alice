package cryptos

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
)

const (
	PubCipherUnknown = PubCipherByte(iota)
	Rsa4096Sha256
	Rsa2048Sha256
	Rsa1024Sha256
)

func NewPubCipher(byte PubCipherByte) IPubCipher {
	switch byte {
	case Rsa4096Sha256:
		return &RsaCipher{size: 4096, hash: NewHash(SHA256)}
	case Rsa2048Sha256:
		return &RsaCipher{size: 2048, hash: NewHash(SHA256)}
	case Rsa1024Sha256:
		return &RsaCipher{size: 1024, hash: NewHash(SHA256)}
	default:
		return nil
	}
}

type RsaCipher struct {
	hash IHash
	size int
}

func (r *RsaCipher) Encrypt(pub []byte, data []byte) ([]byte, error) {
	pubKey, err := parseRsaPubKey(pub)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptOAEP(r.hash.New(), rand.Reader, pubKey, data, nil)
}

func (r *RsaCipher) Decrypt(priv []byte, data []byte) ([]byte, error) {
	privKey, err := parseRsaPrivKey(priv)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptOAEP(r.hash.New(), rand.Reader, privKey, data, nil)
}

func (r *RsaCipher) Sign(priv []byte, data []byte) ([]byte, error) {
	privKey, err := parseRsaPrivKey(priv)
	if err != nil {
		return nil, err
	}

	digest := r.hash.HashSum(data)
	return rsa.SignPSS(rand.Reader, privKey, HashCrypto(r.hash.Byte()), digest, nil)
}

func (r *RsaCipher) Verify(pub []byte, data []byte, sig []byte) error {
	pubKey, err := parseRsaPubKey(pub)
	if err != nil {
		return err
	}

	digest := r.hash.HashSum(data)
	err = rsa.VerifyPSS(pubKey, HashCrypto(r.hash.Byte()), digest, sig, nil)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrBrokenSignature, err.Error())
	}

	return nil
}

func (r *RsaCipher) GeneratePair() (priv, pub []byte, err error) {
	privKey, err := rsa.GenerateKey(rand.Reader, r.size)
	if err != nil {
		return nil, nil, err
	}

	return marshallRsaPrivKey(privKey),
		marshallRsaPubKey(&privKey.PublicKey),
		nil
}

func marshallRsaPrivKey(key *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(key)
}

func parseRsaPrivKey(data []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(data)
}

func marshallRsaPubKey(key *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(key)
}

func parseRsaPubKey(data []byte) (*rsa.PublicKey, error) {
	return x509.ParsePKCS1PublicKey(data)
}
