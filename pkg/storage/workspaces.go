package storage

import (
	"context"
	"fmt"

	"github.com/oka-is/alice/pkg/domain"
)

func (s *Storage) CreateWorkspace(ctx context.Context, uw *domain.UserWorkspace, workspace *domain.Workspace) error {
	return s.Tx(ctx, nil, func(c context.Context, tx *Tx) error {
		return s.createWorkspace(c, tx, uw, workspace)
	})
}

func (s *Storage) DeleteWorkspace(ctx context.Context, ID string) error {
	query := Builder().Delete("workspaces CASCADE").Where("id = ?", ID)
	_, err := s.Exec(ctx, s.db, query)
	return err
}

func (s *Storage) createWorkspace(ctx context.Context, db IConn, uw *domain.UserWorkspace, workspace *domain.Workspace) error {
	err := s.insertWorkspace(ctx, db, workspace)
	if err != nil {
		return fmt.Errorf("failed to insert workspace: %w", err)
	}

	uw.WorkspaceID = workspace.ID
	err = s.insertUserWorkspace(ctx, db, uw)
	if err != nil {
		return fmt.Errorf("failed to insert user workspace: %w", err)
	}

	return nil
}

func (s *Storage) insertWorkspace(ctx context.Context, conn IConn, workspace *domain.Workspace) error {
	query := Builder().
		Insert("workspaces").
		Columns("title_enc").
		Values(workspace.TitleEnc).
		Suffix("RETURNING id, created_at")

	return s.QueryRow(ctx, conn, query).Scan(&workspace.ID, &workspace.CreatedAt)
}
