package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/domain"
)

func TestValidator_ValidateTerminate(t *testing.T) {
	type args struct {
		desc    string
		build   func() ValidateTerminateOpts
		wantErr string
	}

	tests := []args{
		{
			desc: "when valid",
			build: func() ValidateTerminateOpts {
				opts := mustBuildValidateTerminate(t)
				return opts
			},
		}, {
			desc: "when identity mismatch",
			build: func() ValidateTerminateOpts {
				opts := mustBuildValidateTerminate(t)
				opts.Identity = "foo"
				return opts
			},
			wantErr: "Identity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := New().ValidateTerminate(tt.build())

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func mustBuildValidateTerminate(t *testing.T) ValidateTerminateOpts {
	identity := "Identity"

	return ValidateTerminateOpts{
		Identity: identity,
		User: domain.User{
			Identity: domain.NewEmptyString(identity),
		},
	}
}
