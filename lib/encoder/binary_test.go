package encoder

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Units(t *testing.T) {
	type args struct {
		input int
		want  int
	}

	tests := []args{
		{
			input: Int8Size,
			want:  1,
		}, {
			input: Int16Size,
			want:  2,
		}, {
			input: Int32Size,
			want:  4,
		}, {
			input: Int64Size,
			want:  8,
		}, {
			input: Uint8Size,
			want:  1,
		}, {
			input: Uint16Size,
			want:  2,
		}, {
			input: Uint32Size,
			want:  4,
		}, {
			input: Uint64Size,
			want:  8,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input <%d>", tt.input), func(t *testing.T) {
			require.Equal(t, tt.input, tt.want)
		})
	}
}
