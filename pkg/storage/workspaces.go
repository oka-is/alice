package storage

import (
	"context"
	"fmt"

	"github.com/wault-pw/alice/pkg/domain"
)

func (s *Storage) CreateWorkspace(ctx context.Context, uw *domain.UserWorkspace, workspace *domain.Workspace) error {
	return s.Tx(ctx, nil, func(c context.Context, tx IConn) error {
		return s.createWorkspace(c, tx, uw, workspace)
	})
}

func (s *Storage) DeleteWorkspace(ctx context.Context, ID string) error {
	query := Builder().Delete("workspaces CASCADE").Where("id = ?", ID)
	_, err := s.Exec(ctx, s.db, query)
	return err
}

func (s *Storage) FindWorkspace(ctx context.Context, ID string) (out domain.Workspace, err error) {
	query := Builder().Select("*").From("workspaces").Where("id = ?", ID).Limit(1)
	err = s.Get(ctx, s.db, &out, query)
	return
}

func (s *Storage) UpdateWorkspace(ctx context.Context, ID string, titleEnc []byte) error {
	return s.updateWorkspace(ctx, s.db, ID, titleEnc)
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

func (s *Storage) updateWorkspace(ctx context.Context, db IConn, ID string, titleEnc []byte) error {
	query := Builder().Update("workspaces").Set("title_enc", domain.NewEmptyBytes(titleEnc)).Where("id = ? ", ID)
	return s.Exec1(ctx, db, query)
}

func (s *Storage) insertWorkspace(ctx context.Context, conn IConn, workspace *domain.Workspace) error {
	query := Builder().
		Insert("workspaces").
		Columns("title_enc").
		Values(workspace.TitleEnc).
		Suffix("RETURNING id, created_at")

	return s.QueryRow(ctx, conn, query).Scan(&workspace.ID, &workspace.CreatedAt)
}
