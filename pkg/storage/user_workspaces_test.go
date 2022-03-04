package storage

import (
	"context"
	"testing"

	"github.com/oka-is/alice/pkg/domain"
)

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
