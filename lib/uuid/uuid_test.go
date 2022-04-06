package uuid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafe(t *testing.T) {
	type args struct {
		desc  string
		input string
		want  string
	}

	tests := []args{
		{
			desc:  "when ok",
			input: "6aa4c2f4-a77c-4a0f-8278-90a6d6f04eb6",
			want:  "6aa4c2f4-a77c-4a0f-8278-90a6d6f04eb6",
		}, {
			desc:  "when with whitespaces",
			input: "    6aa4c2f4-a77c-4a0f-8278-90a6d6f04eb6   ",
			want:  "6aa4c2f4-a77c-4a0f-8278-90a6d6f04eb6",
		}, {
			desc:  "when unknown",
			input: "foo",
			want:  "00000000-0000-0000-0000-000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			require.Equal(t, tt.want, Safe(tt.input))
		})
	}
}
