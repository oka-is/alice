package storage

import (
	"context"

	"github.com/wault-pw/alice/pkg/domain"
)

func (s *Storage) FindUserWorkspace(ctx context.Context, ID string) (out domain.UserWorkspace, err error) {
	query := Builder().Select("*").From("user_workspaces").Where("id = ?", ID).Limit(1)
	err = s.Get(ctx, s.db, &out, query)
	return
}

func (s *Storage) insertUserWorkspace(ctx context.Context, conn IConn, userWorkspace *domain.UserWorkspace) error {
	query := Builder().
		Insert("user_workspaces").
		Columns(
			"user_id",
			"owner_id",
			"workspace_id",
			"aed_key_enc").
		Values(
			userWorkspace.UserID,
			userWorkspace.OwnerID,
			userWorkspace.WorkspaceID,
			userWorkspace.AedKeyEnc).
		Suffix("RETURNING id, created_at")

	return s.QueryRow(ctx, conn, query).Scan(&userWorkspace.ID, &userWorkspace.CreatedAt)
}
