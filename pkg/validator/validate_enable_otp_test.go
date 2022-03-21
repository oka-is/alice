package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/domain"
)

func TestValidator_ValidateEnableUserOtp(t *testing.T) {
	type args struct {
		desc    string
		build   func() ValidateEnableUserOtpOpts
		wantErr string
	}

	tests := []args{
		{
			desc: "is valid",
			build: func() ValidateEnableUserOtpOpts {
				opts := mustBuildValidateEnableUserOtpOpts(t)
				return opts
			},
		}, {
			desc: "when identity mismatch",
			build: func() ValidateEnableUserOtpOpts {
				opts := mustBuildValidateEnableUserOtpOpts(t)
				opts.Identity = "0"
				return opts
			},
			wantErr: "Identity",
		}, {
			desc: "when secret mismatch",
			build: func() ValidateEnableUserOtpOpts {
				opts := mustBuildValidateEnableUserOtpOpts(t)
				opts.Secret = []byte{0}
				return opts
			},
			wantErr: "OtpSecret",
		}, {
			desc: "when secret is currently set",
			build: func() ValidateEnableUserOtpOpts {
				opts := mustBuildValidateEnableUserOtpOpts(t)
				opts.User.OtpSecret = domain.NewEmptyBytes([]byte{1})
				return opts
			},
			wantErr: "OtpSecret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := New().ValidateEnableUserOtp(tt.build())

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func mustBuildValidateEnableUserOtpOpts(t *testing.T) ValidateEnableUserOtpOpts {
	identity := "foo"
	secret := []byte{1}

	return ValidateEnableUserOtpOpts{
		Identity: identity,
		Secret:   secret,
		User: domain.User{
			OtpSecret:    domain.NewNullBytea(),
			Identity:     domain.NewEmptyString(identity),
			OtpCandidate: domain.NewEmptyBytes(secret),
		},
	}
}
