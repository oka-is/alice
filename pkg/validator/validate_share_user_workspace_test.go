package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/domain"
)

func TestValidator_ValidateShareUserWorkspace(t *testing.T) {
	type args struct {
		desc    string
		build   func() ValidateShareUserWorkspaceOpts
		wantErr string
	}

	tests := []args{
		{
			desc: "is valid",
			build: func() ValidateShareUserWorkspaceOpts {
				opts := mustBuildValidateShareUserWorkspaceOpts(t)
				return opts
			},
		}, {
			desc: "when currently shared",
			build: func() ValidateShareUserWorkspaceOpts {
				opts := mustBuildValidateShareUserWorkspaceOpts(t)
				opts.NotShared = false
				return opts
			},
			wantErr: "User",
		}, {
			desc: "when AedKeyEnc blank",
			build: func() ValidateShareUserWorkspaceOpts {
				opts := mustBuildValidateShareUserWorkspaceOpts(t)
				opts.UserWorkspace.AedKeyEnc = domain.NewNullBytea()
				return opts
			},
			wantErr: "AedKeyEnc",
		}, {
			desc: "when share with self",
			build: func() ValidateShareUserWorkspaceOpts {
				opts := mustBuildValidateShareUserWorkspaceOpts(t)
				opts.UserWorkspace.UserID = opts.UserWorkspace.OwnerID
				return opts
			},
			wantErr: "User",
		}, {
			desc: "when user not exists",
			build: func() ValidateShareUserWorkspaceOpts {
				opts := mustBuildValidateShareUserWorkspaceOpts(t)
				opts.UserExists = false
				return opts
			},
			wantErr: "User",
		}, {
			desc: "when sharing restricted by a user",
			build: func() ValidateShareUserWorkspaceOpts {
				opts := mustBuildValidateShareUserWorkspaceOpts(t)
				opts.SharingPermitted = false
				return opts
			},
			wantErr: "User",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := New().ValidateShareUserWorkspace(tt.build())

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func mustBuildValidateShareUserWorkspaceOpts(t *testing.T) ValidateShareUserWorkspaceOpts {
	return ValidateShareUserWorkspaceOpts{
		NotShared:        true,
		UserExists:       true,
		SharingPermitted: true,
		UserWorkspace: domain.UserWorkspace{
			AedKeyEnc: domain.NewEmptyBytes([]byte{1}),
			UserID:    domain.NewEmptyString("workspace-user"),
			OwnerID:   domain.NewEmptyString("workspace-owner"),
		},
	}
}
