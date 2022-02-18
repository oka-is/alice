package null

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parsePgArrayHex(t *testing.T) {
	type args struct {
		desc  string
		input string
		want  [][]byte
	}
	tests := []args{
		{
			desc:  "when empty",
			input: "",
			want:  [][]byte{},
		}, {
			desc:  "when one",
			input: `{"\\x756e6f"}`,
			want: [][]byte{
				[]byte("uno"),
			},
		}, {
			desc:  "when two",
			input: `{"\\x756e6f","\\x646f73"}`,
			want: [][]byte{
				[]byte("uno"),
				[]byte("dos"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := parsePgArrayHex(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
