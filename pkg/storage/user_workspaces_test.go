package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/domain"
)

func TestStorage_FindUserWorkspaceLink(t *testing.T) {
	store, savepoint := MustNewStore(t)
	t.Cleanup(savepoint.Flush)

	uw1 := mustCreateUserWorkspace(t, store, &domain.UserWorkspace{})
	uw2 := mustCreateUserWorkspace(t, store, &domain.UserWorkspace{})

	type args struct {
		desc        string
		userID      string
		workspaceID string
		wantID      string
		wantErr     error
	}

	tests := []args{
		{
			desc:        "when ok",
			userID:      uw1.UserID.String,
			workspaceID: uw1.WorkspaceID.String,
			wantID:      uw1.ID.String,
		}, {
			desc:        "error when broken userID",
			userID:      uw2.UserID.String,
			workspaceID: uw1.WorkspaceID.String,
			wantErr:     ErrNotFound,
		}, {
			desc:        "error when broken workspaceID",
			userID:      uw1.UserID.String,
			workspaceID: uw2.WorkspaceID.String,
			wantErr:     ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ctx := context.Background()
			got, err := store.FindUserWorkspaceLink(ctx, tt.userID, tt.workspaceID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantID, got.ID.String)
		})
	}
}

func mustBuildUserWorkspace(t *testing.T, storage *Storage, input *domain.UserWorkspace) *domain.UserWorkspace {
	out := &domain.UserWorkspace{
		AedKeyEnc: domain.NewEmptyBytes([]byte("AedKeyEnc")),
	}

	if input.UserID.Valid {
		out.UserID = input.UserID
	} else {
		out.UserID = mustCreateUser(t, storage, &domain.User{}).ID
	}

	if input.OwnerID.Valid {
		out.OwnerID = input.OwnerID
	} else {
		out.OwnerID = mustCreateUser(t, storage, &domain.User{}).ID
	}

	if input.WorkspaceID.Valid {
		out.WorkspaceID = input.WorkspaceID
	} else {
		out.WorkspaceID = mustCreateWorkspace(t, storage, &domain.Workspace{}).ID
	}

	if input.AedKeyEnc.Valid {
		out.AedKeyEnc = input.AedKeyEnc
	}

	return out
}

func mustCreateUserWorkspace(t *testing.T, storage *Storage, input *domain.UserWorkspace) *domain.UserWorkspace {
	ctx := context.Background()
	output := mustBuildUserWorkspace(t, storage, input)
	if err := storage.insertUserWorkspace(ctx, storage.db, output); err != nil {
		t.Errorf("cant create factory user workspace: %s", err)
		t.FailNow()
	}
	return output
}
