package storage

import (
	"context"

	"github.com/oka-is/alice/pkg/domain"
)

func (s *Storage) insertUserWorkspace(ctx context.Context, conn IConn, userWorkspace *domain.UserWorkspace) error {
	query := Builder().
		Insert("user_workspaces").
		Columns(
			"user_id",
			"owner_id",
			"workspace_id",
			"aed_key_enc",
			"aed_key_tag").
		Values(
			userWorkspace.UserID,
			userWorkspace.OwnerID,
			userWorkspace.WorkspaceID,
			userWorkspace.AedKeyEnc,
			userWorkspace.AedKeyTag).
		Suffix("RETURNING id, created_at")

	return s.QueryRow(ctx, conn, query).Scan(&userWorkspace.ID, &userWorkspace.CreatedAt)
}
