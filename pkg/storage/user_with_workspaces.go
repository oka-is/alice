package storage

import (
	"context"

	"github.com/oka-is/alice/pkg/domain"
)

func (s *Storage) ListUserWithWorkspaces(ctx context.Context, userID string) (out []domain.UserWithWorkspace, err error) {
	query := userWorkspacesScope().Where("user_workspaces.user_id = ?", userID)
	err = s.Select(ctx, s.db, &out, query)
	return
}

func (s *Storage) FindUserWithWorkspace(ctx context.Context, ID string) (out domain.UserWithWorkspace, err error) {
	query := userWorkspacesScope().Where("user_workspaces.id = ?", ID).Limit(1)
	err = s.Get(ctx, s.db, &out, query)
	return
}

func userWorkspacesScope() SelectBuilder {
	return Builder().
		Select().
		Columns(
			"user_workspaces.id AS record_id",
			"user_workspaces.user_id AS user_id",
			"user_workspaces.owner_id AS owner_id",
			"owners.pub_key AS owner_pub_key",
			"user_workspaces.workspace_id AS workspace_id",
			"user_workspaces.aed_key_enc AS aed_key_enc",
			"workspaces.title_enc AS title_enc",
			"user_workspaces.created_at AS record_created_at",
			"workspaces.created_at AS workspace_created_at").
		From("user_workspaces").
		InnerJoin("workspaces ON workspaces.id = user_workspaces.workspace_id").
		InnerJoin("users owners ON owners.id = user_workspaces.owner_id")
}
