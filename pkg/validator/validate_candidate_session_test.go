package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidator_ValidateCandidateSession(t *testing.T) {
	type args struct {
		desc    string
		build   func() ValidateCandidateSessionOpts
		wantErr string
	}

	tests := []args{
		{
			desc: "is valid",
			build: func() ValidateCandidateSessionOpts {
				opts := mustBuildValidateCandidateSessionOpts(t)
				return opts
			},
		}, {
			desc: "valid when attempts is less then limit",
			build: func() ValidateCandidateSessionOpts {
				opts := mustBuildValidateCandidateSessionOpts(t)
				opts.Attempts -= 1
				return opts
			},
		}, {
			desc: "invalid when attempts is greater then limit",
			build: func() ValidateCandidateSessionOpts {
				opts := mustBuildValidateCandidateSessionOpts(t)
				opts.Attempts += 1
				return opts
			},
			wantErr: "Attempts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := New().ValidateCandidateSession(tt.build())

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func mustBuildValidateCandidateSessionOpts(t *testing.T) ValidateCandidateSessionOpts {
	return ValidateCandidateSessionOpts{
		Attempts: MaxLoginAttempts,
	}
}
