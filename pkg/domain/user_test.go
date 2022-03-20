package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUser_IsOtpEnabled(t *testing.T) {
	type args struct {
		desc  string
		input User
		want  bool
	}

	tests := []args{
		{
			desc:  "when false",
			input: User{},
			want:  false,
		}, {
			desc:  "when true",
			input: User{OtpSecret: NewEmptyBytes([]byte{1})},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.input.IsOtpEnabled()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUser_IsOtpDisabled(t *testing.T) {
	type args struct {
		desc  string
		input User
		want  bool
	}

	tests := []args{
		{
			desc:  "when true",
			input: User{},
			want:  true,
		}, {
			desc:  "when false",
			input: User{OtpSecret: NewEmptyBytes([]byte{1})},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.input.IsOtpDisabled()
			require.Equal(t, tt.want, got)
		})
	}
}
