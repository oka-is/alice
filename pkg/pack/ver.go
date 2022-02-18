package pack

import (
	"github.com/oka-is/alice/lib/cryptos"
	"github.com/oka-is/srp6ago"
)

type VerByte byte

const (
	VerUnknown = VerByte(iota)
	Ver1
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
)

func NewWer(byte VerByte) *Ver {
	switch byte {
	case Ver1:
		return ver1
	default:
		return nil
	}
}

func (v *Ver) NewSrpServer(verifier, salt []byte) *srp6ago.Server {
	return srp6ago.NewServer(verifier, salt, v.SrpGroup)
}

//func (v *Ver) Derive() cryptos.IDerive {
//	return cryptos.NewDerive(v.DeriveByte)
//}
//
//func (v *Ver) DeriveSaltSize() int {
//	return cryptos.DeriveSaltSize(v.DeriveByte)
//}
//
//func (v *Ver) AedCipherSizes() (keySize, ivSize, blockSize int) {
//	return cryptos.AedCipherSizes(v.AedCipherByte)
//}
//
//func (v *Ver) AedKeySize() int {
//	key, _, _ := v.AedCipherSizes()
//	return key
//}
//
//func (v *Ver) PubCipher() cryptos.IPubCipher {
//	return cryptos.NewPubCipher(v.PubCipherByte)
//}
//
//func (v *Ver) AedIVSize() int {
//	_, iv, _ := v.AedCipherSizes()
//	return iv
//}
//
//func (v *Ver) MacSize() int {
//	return cryptos.MacSize(v.MacByte)
//}
//
//func (v *Ver) AedCipher() cryptos.IAedCipher {
//	return cryptos.NewAedCipher(v.AedCipherByte)
//}
//
//func (v *Ver) Mac() cryptos.IMac {
//	return cryptos.NewMac(v.MacByte)
//}
//
//func (v *Ver) RandomIV() ([]byte, error) {
//	return cryptos.SecureRandSize(v.AedIVSize())
//}
//
//func (v *Ver) PubEncrypt(pubKey, data []byte) ([]byte, error) {
//	return v.PubCipher().Encrypt(pubKey, data)
//}
//
//func (v *Ver) PrivDecrypt(privKey, data []byte) ([]byte, error) {
//	return v.PubCipher().Decrypt(privKey, data)
//}
//
//func (v *Ver) PrivSign(privKey, data []byte) ([]byte, error) {
//	return v.PubCipher().Sign(privKey, data)
//}
//
//func (v *Ver) PubVerify(pubKey, data, sig []byte) error {
//	return v.PubCipher().Verify(pubKey, data, sig)
//}
//
//// AedEncrypt encrypts plaintext and signs plaintext+addon
//func (v *Ver) AedEncrypt(key, plaintext, addon []byte) ([]byte, error) {
//	iv, err := v.RandomIV()
//	if err != nil {
//		return nil, fmt.Errorf("failed to generate random iv: %w", err)
//	}
//
//	cipher, err := v.AedCipher().Init(key, iv)
//	if err != nil {
//		return nil, fmt.Errorf("failed to init AED: %w", err)
//	}
//
//	out := cipher.Encrypt(plaintext, addon)
//	return append(iv, out...), nil
//}
//
//// AedDecrypt decrypts ciphertext and verifies ciphertext+addon
//func (v *Ver) AedDecrypt(key, ciphertext, addon []byte) ([]byte, error) {
//	if len(ciphertext) == 0 {
//		return []byte{}, nil
//	}
//
//	is := v.AedIVSize()
//	iv, data := ciphertext[:is], ciphertext[is:]
//
//	cipher, err := v.AedCipher().Init(key, iv)
//	if err != nil {
//		return nil, fmt.Errorf("failed to init AED: %w", err)
//	}
//
//	return cipher.Decrypt(data, addon)
//}
