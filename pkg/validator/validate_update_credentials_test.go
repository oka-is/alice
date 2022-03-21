package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/domain"
)

func TestValidator_ValidateUpdateCredentials(t *testing.T) {
	type args struct {
		desc    string
		build   func() ValidateUpdateCredentialsOpts
		wantErr string
	}

	tests := []args{
		{
			desc: "is valid",
			build: func() ValidateUpdateCredentialsOpts {
				opts := mustBuildValidateUpdateCredentialsOpts(t)
				return opts
			},
		}, {
			desc: "when Identity mismatch",
			build: func() ValidateUpdateCredentialsOpts {
				opts := mustBuildValidateUpdateCredentialsOpts(t)
				opts.OldIdentity = "foo"
				return opts
			},
			wantErr: "OldIdentity",
		}, {
			desc: "when Identity blank",
			build: func() ValidateUpdateCredentialsOpts {
				opts := mustBuildValidateUpdateCredentialsOpts(t)
				opts.NewUser.Identity = domain.NewNullString()
				return opts
			},
			wantErr: "Identity",
		}, {
			desc: "when Verifier blank",
			build: func() ValidateUpdateCredentialsOpts {
				opts := mustBuildValidateUpdateCredentialsOpts(t)
				opts.NewUser.Verifier = domain.NewNullBytea()
				return opts
			},
			wantErr: "Verifier",
		}, {
			desc: "when SrpSalt blank",
			build: func() ValidateUpdateCredentialsOpts {
				opts := mustBuildValidateUpdateCredentialsOpts(t)
				opts.NewUser.SrpSalt = domain.NewNullBytea()
				return opts
			},
			wantErr: "SrpSalt",
		}, {
			desc: "when PasswdSalt blank",
			build: func() ValidateUpdateCredentialsOpts {
				opts := mustBuildValidateUpdateCredentialsOpts(t)
				opts.NewUser.PasswdSalt = domain.NewNullBytea()
				return opts
			},
			wantErr: "PasswdSalt",
		}, {
			desc: "when PrivKeyEnc blank",
			build: func() ValidateUpdateCredentialsOpts {
				opts := mustBuildValidateUpdateCredentialsOpts(t)
				opts.NewUser.PrivKeyEnc = domain.NewNullBytea()
				return opts
			},
			wantErr: "PrivKeyEnc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := New().ValidateUpdateCredentials(tt.build())

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func mustBuildValidateUpdateCredentialsOpts(t *testing.T) ValidateUpdateCredentialsOpts {
	oldIdentity := "oldIdentity"

	return ValidateUpdateCredentialsOpts{
		OldIdentity: oldIdentity,
		OldUser: domain.User{
			Identity: domain.NewEmptyString(oldIdentity),
		},
		NewUser: domain.User{
			Identity:   domain.NewEmptyString("Identity"),
			Verifier:   domain.NewEmptyBytes([]byte{1}),
			SrpSalt:    domain.NewEmptyBytes([]byte{1}),
			PasswdSalt: domain.NewEmptyBytes([]byte{1}),
			PrivKeyEnc: domain.NewEmptyBytes([]byte{1}),
		},
	}
}
