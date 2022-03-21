package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/domain"
)

func TestValidator_ValidateUser(t *testing.T) {
	type args struct {
		desc    string
		build   func() domain.User
		wantErr string
	}

	tests := []args{
		{
			desc: "is valid",
			build: func() domain.User {
				opts := mustBuildValidateUser(t)
				return opts
			},
		}, {
			desc: "when Ver blank",
			build: func() domain.User {
				opts := mustBuildValidateUser(t)
				opts.Ver = domain.NewEmptyInt64(0)
				return opts
			},
			wantErr: "Ver",
		}, {
			desc: "when Identity blank",
			build: func() domain.User {
				opts := mustBuildValidateUser(t)
				opts.Identity = domain.NewNullString()
				return opts
			},
			wantErr: "Identity",
		}, {
			desc: "when Verifier blank",
			build: func() domain.User {
				opts := mustBuildValidateUser(t)
				opts.Verifier = domain.NewNullBytea()
				return opts
			},
			wantErr: "Verifier",
		}, {
			desc: "when SrpSalt blank",
			build: func() domain.User {
				opts := mustBuildValidateUser(t)
				opts.SrpSalt = domain.NewNullBytea()
				return opts
			},
			wantErr: "SrpSalt",
		}, {
			desc: "when PasswdSalt blank",
			build: func() domain.User {
				opts := mustBuildValidateUser(t)
				opts.PasswdSalt = domain.NewNullBytea()
				return opts
			},
			wantErr: "PasswdSalt",
		}, {
			desc: "when PrivKeyEnc blank",
			build: func() domain.User {
				opts := mustBuildValidateUser(t)
				opts.PrivKeyEnc = domain.NewNullBytea()
				return opts
			},
			wantErr: "PrivKeyEnc",
		}, {
			desc: "when PubKey blank",
			build: func() domain.User {
				opts := mustBuildValidateUser(t)
				opts.PubKey = domain.NewNullBytea()
				return opts
			},
			wantErr: "PubKey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := New().ValidateUser(tt.build())

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func mustBuildValidateUser(t *testing.T) domain.User {
	return domain.User{
		Ver:        domain.NewNullInt64(1),
		Identity:   domain.NewEmptyString("Identity"),
		Verifier:   domain.NewEmptyBytes([]byte{1}),
		SrpSalt:    domain.NewEmptyBytes([]byte{1}),
		PasswdSalt: domain.NewEmptyBytes([]byte{1}),
		PrivKeyEnc: domain.NewEmptyBytes([]byte{1}),
		PubKey:     domain.NewEmptyBytes([]byte{1}),
	}
}
