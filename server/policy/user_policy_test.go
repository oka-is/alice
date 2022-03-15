package policy

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/domain"
)

func TestUserPolicy_CanWrite(t *testing.T) {
	type args struct {
		desc string
		user domain.User
		want error
	}

	tests := []args{
		{
			desc: "when readonly",
			user: domain.User{Readonly: domain.NewEmptyBool(true)},
			want: ErrDenied,
		}, {
			desc: "when not readonly",
			user: domain.User{Readonly: domain.NewEmptyBool(false)},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			policy := new(UserPolicy)
			policy.Wrap(tt.user)
			require.Equal(t, tt.want, policy.CanWrite())
		})
	}
}
