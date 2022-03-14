package cryptos

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/lib/encoder"
)

// @refs https://emn178.github.io/online-tools/sha256.html
func TestHash_New(t *testing.T) {
	type args struct {
		desc  string
		byte  HashByte
		input string
		want  string
	}

	tests := []args{
		{
			desc:  "it works for SHA256",
			byte:  SHA256,
			input: "foo",
			want:  "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae",
		}, {
			desc:  "it works for SHA512",
			byte:  SHA512,
			input: "foo",
			want:  "f7fbba6e0636f890e56fbbf3283e524c6fa3204ae298382d624741d0dc6638326e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := NewHash(tt.byte).HashSum([]byte(tt.input))
			require.Equal(t, tt.want, encoder.B2Hex(got))
		})
	}
}
