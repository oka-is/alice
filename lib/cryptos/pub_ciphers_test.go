package cryptos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRsaCipher(t *testing.T) {
	t.Run("it works", func(t *testing.T) {
		enc, dec := NewPubCipher(Rsa2048Sha256), NewPubCipher(Rsa2048Sha256)
		plaintext := []byte{1, 2, 3}

		priv, pub, err := enc.GeneratePair()
		require.NoError(t, err)

		ciphertext, err := enc.Encrypt(pub, plaintext)
		require.NoError(t, err)

		tag, err := enc.Sign(priv, ciphertext)
		require.NoError(t, err)

		got, err := dec.Decrypt(priv, ciphertext)
		require.NoError(t, err)
		require.Equal(t, plaintext, got)

		err = dec.Verify(pub, ciphertext, tag)
		require.NoError(t, err)
	})
}
