package cryptos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHmacSha256(t *testing.T) {
	key := []byte{1, 2, 3, 4, 5}
	input := []byte{99, 98, 97, 96}
	want := []byte{0x1d, 0xb, 0x54, 0x8d, 0x99, 0xa7, 0x96, 0x4,
		0xa6, 0x8a, 0xf2, 0xef, 0x23, 0x87, 0xf9, 0x60, 0xec,
		0x11, 0xa3, 0x60, 0x97, 0x77, 0xc7, 0x1c, 0xed, 0x5b,
		0x76, 0x74, 0xd4, 0xf1, 0x52, 0xa1}

	t.Run("it works for whole peace", func(t *testing.T) {
		mac := NewMac(HmacSha256).Init(key)
		_, err := mac.Write(input)
		require.NoError(t, err)
		require.Equal(t, want, mac.Sum())
	})

	t.Run("it works for small peaces", func(t *testing.T) {
		mac := NewMac(HmacSha256).Init(key)

		for _, b := range input {
			_, err := mac.Write([]byte{b})
			require.NoError(t, err)
		}

		require.Equal(t, want, mac.Sum())
	})
}
